package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type PaymentRepositoryStr struct {
	db.ModeledDatabase[model.Payment]
}

func NewPaymentRepository(service db.DatabaseService) PaymentRepository {
	return &PaymentRepositoryStr{db.NewModeledDatabase(model.Payment{}, service)}
}

func (p *PaymentRepositoryStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(db.Database{}),
	}
}

func (p *PaymentRepositoryStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewPaymentRepository(dependency.FindRequiredDependency[db.Database, db.DatabaseService](factory))
}

type PaymentRepository interface {
	db.ModeledDatabase[model.Payment]
	transactional.Transactional[PaymentRepository]
	FindByHouseId(houseId uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto
	FindByUserId(userId uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto
	FindByGroupId(id uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto
	FindByProviderId(providerId uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto
}

func (p *PaymentRepositoryStr) FindByHouseId(houseId uuid.UUID, limit int, offset int, from, to *time.Time) (response []model.PaymentDto) {
	whereQuery := "house_id = ?"
	whereArgs := []any{houseId}

	if from != nil && to != nil {
		whereQuery += " AND date BETWEEN ? AND ?"
		whereArgs = append(whereArgs, from, to)
	} else if from != nil {
		whereQuery += " AND date >= ?"
		whereArgs = append(whereArgs, from)
	}

	query := p.Modeled().
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

func (p *PaymentRepositoryStr) FindByUserId(userId uuid.UUID, limit int, offset int, from, to *time.Time) (response []model.PaymentDto) {
	whereQuery := "user_id = ?"
	whereArgs := []any{userId}

	if from != nil && to != nil {
		whereQuery += " AND date BETWEEN ? AND ?"
		whereArgs = append(whereArgs, from, to)
	} else if from != nil {
		whereQuery += " AND date >= ?"
		whereArgs = append(whereArgs, from)
	}

	query := p.Modeled().
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

func (p *PaymentRepositoryStr) FindByProviderId(providerId uuid.UUID, limit int, offset int, from, to *time.Time) (response []model.PaymentDto) {
	whereQuery := "provider_id = ?"
	whereArgs := []any{providerId}

	if from != nil && to != nil {
		whereQuery += " AND date BETWEEN ? AND ?"
		whereArgs = append(whereArgs, from, to)
	} else if from != nil {
		whereQuery += " AND date >= ?"
		whereArgs = append(whereArgs, from)
	}

	err := p.Modeled().
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

func (p *PaymentRepositoryStr) FindByGroupId(id uuid.UUID, limit int, offset int, from, to *time.Time) (response []model.PaymentDto) {
	whereQuery := "hg.group_id = ?"
	whereArgs := []any{id}

	if from != nil && to != nil {
		whereQuery += " AND date BETWEEN ? AND ?"
		whereArgs = append(whereArgs, from, to)
	} else if from != nil {
		whereQuery += " AND date >= ?"
		whereArgs = append(whereArgs, from)
	}

	if err := p.Modeled().
		Joins("full join houses h ON h.id = payments.house_id full join house_groups hg ON hg.house_id = h.id").
		Where(whereQuery, whereArgs...).
		Order("date desc").
		Limit(limit).
		Offset(offset).
		Find(&response).Error; err != nil {
		log.Error().Err(err)
	}

	return response
}

func (p *PaymentRepositoryStr) Transactional(tx *gorm.DB) PaymentRepository {
	return &PaymentRepositoryStr{
		ModeledDatabase: db.NewTransactionalModeledDatabase(p.GetEntity(), tx),
	}
}
