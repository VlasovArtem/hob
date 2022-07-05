package repository

import (
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/db"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type IncomeRepositoryStr struct {
	db.ModeledDatabase[model.Income]
}

func (i *IncomeRepositoryStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewIncomeRepository(factory.FindRequiredByObject(db.Database{}).(db.DatabaseService))
}

func NewIncomeRepository(database db.DatabaseService) IncomeRepository {
	return &IncomeRepositoryStr{db.NewModeledDatabase(model.Income{}, database)}
}

type IncomeRepository interface {
	db.ModeledDatabase[model.Income]
	transactional.Transactional[IncomeRepository]
	FindById(id uuid.UUID) (response model.Income, err error)
	FindByHouseId(id uuid.UUID, limit int, offset int, from, to *time.Time) ([]model.IncomeDto, error)
	FindByGroupIds(groupIds []uuid.UUID, limit int, offset int, from, to *time.Time) ([]model.IncomeDto, error)
	UpdateByRequest(id uuid.UUID, request model.UpdateIncomeRequest) error
}

func (i *IncomeRepositoryStr) FindById(id uuid.UUID) (response model.Income, err error) {
	response.Id = id
	if err = i.DB().Preload("Groups").First(&response).Error; err != nil {
		return model.Income{}, err
	}
	return response, err
}

func (i *IncomeRepositoryStr) FindByHouseId(id uuid.UUID, limit int, offset int, from, to *time.Time) (response []model.IncomeDto, err error) {
	var responseEntities []model.Income

	whereQuery := "(incomes.house_id = ? OR hg.house_id = ?)"
	whereArgs := []any{id, id}

	if from != nil && to != nil {
		whereQuery += " AND date BETWEEN ? AND ?"
		whereArgs = append(whereArgs, from, to)
	} else if from != nil {
		whereQuery += " AND date >= ?"
		whereArgs = append(whereArgs, from)
	}

	query := i.DB().
		Joins("FULL JOIN income_groups ig ON ig.income_id = incomes.id FULL JOIN house_groups hg ON hg.group_id = ig.group_id").
		Order("incomes.date desc").
		Where(whereQuery, whereArgs...).
		Preload("Groups")

	if limit >= 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.
		Limit(limit).
		Offset(offset).
		Find(&responseEntities).Error; err != nil {
		return []model.IncomeDto{}, err
	}

	return common.MapSlice(responseEntities, func(i model.Income) model.IncomeDto {
		return i.ToDto()
	}), nil
}

func (i *IncomeRepositoryStr) FindByGroupIds(groupIds []uuid.UUID, limit int, offset int, from, to *time.Time) (response []model.IncomeDto, err error) {
	var responseEntity []model.Income

	whereQuery := ""
	whereArgs := []any{}

	if from != nil && to != nil {
		whereQuery += "date BETWEEN ? AND ?"
		whereArgs = append(whereArgs, from, to)
	} else if from != nil {
		if len(whereArgs) == 0 {
			whereQuery += "date >= ?"
		} else {
			whereQuery += " AND date >= ?"
		}
		whereArgs = append(whereArgs, from)
	}

	query := i.DB().
		Order("incomes.date desc").
		Where(whereQuery, whereArgs...).
		Joins("JOIN income_groups ON income_groups.income_id = incomes.id AND income_groups.group_id IN ?", groupIds).
		Preload("Groups")

	if limit >= 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err = query.Find(&responseEntity).Error; err != nil {
		return []model.IncomeDto{}, err
	}

	return common.MapSlice(responseEntity, func(entity model.Income) model.IncomeDto {
		return entity.ToDto()
	}), nil
}

func (i *IncomeRepositoryStr) UpdateByRequest(id uuid.UUID, request model.UpdateIncomeRequest) error {
	err := i.Update(id, struct {
		Name        string
		Description string
		Date        time.Time
		Sum         float32
	}{
		request.Name,
		request.Description,
		request.Date,
		request.Sum,
	})

	if err != nil {
		return err
	}

	entity, err := i.FindById(id)

	if err != nil {
		return err
	}

	var groups = common.MapSlice(request.GroupIds, groupModel.GroupIdToGroup)

	return i.DBModeled(&entity).Association("Groups").Replace(groups)
}

func (i *IncomeRepositoryStr) CalculateSumByHouseId(houseId uuid.UUID, from *time.Time, sum *float64) error {
	if from != nil {
		return i.DB().
			Raw(`SELECT SUM(sum) FROM incomes WHERE house_id = ? AND date > ?`, houseId, from).
			Scan(sum).
			Error
	} else {
		return i.DB().
			Raw(`SELECT SUM(sum) FROM incomes WHERE house_id = ?`, houseId).
			Scan(sum).
			Error
	}
}

func (i *IncomeRepositoryStr) CalculateSumByGroupId(groupId uuid.UUID, from *time.Time, sum *float64) error {
	if from != nil {
		return i.DB().
			Raw(`SELECT SUM(i.sum) FROM incomes i JOIN income_groups ig ON i.id = ig.income_id WHERE ig.group_id = ? AND date > ?`, groupId, from).
			Scan(sum).
			Error
	} else {
		return i.DB().
			Raw(`SELECT SUM(i.sum) FROM incomes i JOIN income_groups ig ON i.id = ig.income_id WHERE ig.group_id = ?`, groupId).
			Scan(sum).
			Error
	}
}

func (i *IncomeRepositoryStr) Transactional(tx *gorm.DB) IncomeRepository {
	return &IncomeRepositoryStr{db.NewTransactionalModeledDatabase(i.GetEntity(), tx)}
}
