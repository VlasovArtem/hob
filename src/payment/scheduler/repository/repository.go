package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type PaymentSchedulerRepositoryStr struct {
	db.ModeledDatabase[model.PaymentScheduler]
}

func NewPaymentSchedulerRepository(database db.DatabaseService) PaymentSchedulerRepository {
	return &PaymentSchedulerRepositoryStr{db.NewModeledDatabase(model.PaymentScheduler{}, database)}
}

func (p *PaymentSchedulerRepositoryStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(db.Database{}),
	}
}

func (p *PaymentSchedulerRepositoryStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewPaymentSchedulerRepository(dependency.FindRequiredDependency[db.Database, db.DatabaseService](factory))
}

type PaymentSchedulerRepository interface {
	transactional.Transactional[PaymentSchedulerRepository]
	db.ModeledDatabase[model.PaymentScheduler]
	FindByHouseId(houseId uuid.UUID) []model.PaymentSchedulerDto
	FindByUserId(userId uuid.UUID) []model.PaymentSchedulerDto
	FindByProviderId(providerId uuid.UUID) []model.PaymentSchedulerDto
}

func (p *PaymentSchedulerRepositoryStr) FindByHouseId(houseId uuid.UUID) (response []model.PaymentSchedulerDto) {
	if err := p.FindReceiverBy(&response, "house_id = ?", houseId); err != nil {
		log.Error().Err(err).Msg("Failed to find payment scheduler by house id")
	}
	return
}

func (p *PaymentSchedulerRepositoryStr) FindByUserId(userId uuid.UUID) (response []model.PaymentSchedulerDto) {
	if err := p.FindReceiverBy(&response, "user_id = ?", userId); err != nil {
		log.Error().Err(err).Msg("Failed to find payment scheduler by house id")
	}
	return
}

func (p *PaymentSchedulerRepositoryStr) FindByProviderId(providerId uuid.UUID) (response []model.PaymentSchedulerDto) {
	if err := p.FindReceiverBy(&response, "provider_id = ?", providerId); err != nil {
		log.Error().Err(err).Msg("Failed to find payment scheduler by house id")
	}
	return
}

func (p *PaymentSchedulerRepositoryStr) Transactional(tx *gorm.DB) PaymentSchedulerRepository {
	return &PaymentSchedulerRepositoryStr{
		ModeledDatabase: db.NewTransactionalModeledDatabase(p.GetEntity(), tx),
	}
}
