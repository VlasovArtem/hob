package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

var entity = model.Payment{}

type PaymentRepositoryObject struct {
	database db.ModeledDatabase
}

func NewPaymentRepository(service db.DatabaseService) PaymentRepository {
	return &PaymentRepositoryObject{
		db.ModeledDatabase{
			DatabaseService: service,
			Model:           entity,
		},
	}
}

func (p *PaymentRepositoryObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewPaymentRepository(dependency.FindRequiredDependency[db.DatabaseObject, db.DatabaseService](factory))
}

func (p *PaymentRepositoryObject) GetEntity() any {
	return entity
}

type PaymentRepository interface {
	Create(entity model.Payment) (model.Payment, error)
	CreateBatch(entities []model.Payment) ([]model.Payment, error)
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) (model.Payment, error)
	FindByHouseId(houseId uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto
	FindByUserId(userId uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto
	FindByProviderId(providerId uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(entity model.Payment) error
	CalculateSum(houseIds []uuid.UUID, from *time.Time, f *float64)
}

func (p *PaymentRepositoryObject) Create(entity model.Payment) (model.Payment, error) {
	return entity, p.database.Create(&entity)
}

func (p *PaymentRepositoryObject) CreateBatch(entities []model.Payment) ([]model.Payment, error) {
	return entities, p.database.Create(&entities)
}

func (p *PaymentRepositoryObject) Delete(id uuid.UUID) (err error) {
	return p.database.Delete(id)
}

func (p *PaymentRepositoryObject) FindById(id uuid.UUID) (response model.Payment, err error) {
	return response, p.database.Find(&response, id)
}

func (p *PaymentRepositoryObject) FindByHouseId(houseId uuid.UUID, limit int, offset int, from, to *time.Time) (response []model.PaymentDto) {
	whereQuery := "house_id = ?"
	whereArgs := []any{houseId}

	if from != nil && to != nil {
		whereQuery += " AND date BETWEEN ? AND ?"
		whereArgs = append(whereArgs, from, to)
	} else if from != nil {
		whereQuery += " AND date >= ?"
		whereArgs = append(whereArgs, from)
	}

	query := p.database.Modeled().
		Where(whereQuery, whereArgs...).
		Order("date desc")

	if limit >= 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.
		Find(&response).
		Error

	if err != nil {
		log.Err(err).Msg("Error during find payments by house id")
		return make([]model.PaymentDto, 0)
	}
	return response
}

func (p *PaymentRepositoryObject) FindByUserId(userId uuid.UUID, limit int, offset int, from, to *time.Time) (response []model.PaymentDto) {
	whereQuery := "user_id = ?"
	whereArgs := []any{userId}

	if from != nil && to != nil {
		whereQuery += " AND date BETWEEN ? AND ?"
		whereArgs = append(whereArgs, from, to)
	} else if from != nil {
		whereQuery += " AND date >= ?"
		whereArgs = append(whereArgs, from)
	}

	query := p.database.Modeled().
		Where(whereQuery, whereArgs...).
		Order("date desc")

	if limit >= 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.
		Find(&response).
		Error

	if err != nil {
		log.Err(err).Msg("Error during find payments by user id")
		return make([]model.PaymentDto, 0)
	}
	return response
}

func (p *PaymentRepositoryObject) FindByProviderId(providerId uuid.UUID, limit int, offset int, from, to *time.Time) (response []model.PaymentDto) {
	whereQuery := "provider_id = ?"
	whereArgs := []any{providerId}

	if from != nil && to != nil {
		whereQuery += " AND date BETWEEN ? AND ?"
		whereArgs = append(whereArgs, from, to)
	} else if from != nil {
		whereQuery += " AND date >= ?"
		whereArgs = append(whereArgs, from)
	}

	err := p.database.Modeled().
		Where(whereQuery, whereArgs...).
		Order("date desc").
		Limit(limit).
		Offset(offset).
		Find(&response).
		Error

	if err != nil {
		log.Err(err).Msg("Error during find payments by provider id")
		return make([]model.PaymentDto, 0)
	}
	return response
}

func (p *PaymentRepositoryObject) findBy(query any, conditions ...any) (response []model.PaymentDto) {
	_ = p.database.FindBy(&response, query, conditions...)

	return response
}

func (p *PaymentRepositoryObject) ExistsById(id uuid.UUID) bool {
	return p.database.Exists(id)
}

func (p *PaymentRepositoryObject) DeleteById(id uuid.UUID) error {
	return p.database.Delete(id)
}

func (p *PaymentRepositoryObject) Update(entity model.Payment) error {
	return p.database.Update(entity.Id, entity, "HouseId", "House", "UserId", "User")
}

func (p *PaymentRepositoryObject) CalculateSum(houseIds []uuid.UUID, from *time.Time, sum *float64) {
	if len(houseIds) == 0 {
		return
	}

	sql := `SELECT SUM(sum) FROM payments WHERE house_id IN (?)%s ORDER BY date DESC`

	if from != nil {
		sql = fmt.Sprintf(sql, " AND date > ?")
	} else {
		sql = fmt.Sprintf(sql, "")
	}

	p.database.D().
		Exec(sql, houseIds, from).
		Scan(sum)
}
