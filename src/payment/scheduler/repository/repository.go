package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type PaymentSchedulerRepositoryObject struct {
	database db.DatabaseService
}

func NewPaymentSchedulerRepository(database db.DatabaseService) PaymentSchedulerRepository {
	return &PaymentSchedulerRepositoryObject{database}
}

func (p *PaymentSchedulerRepositoryObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return NewPaymentSchedulerRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (p *PaymentSchedulerRepositoryObject) GetEntity() interface{} {
	return model.PaymentScheduler{}
}

type PaymentSchedulerRepository interface {
	Create(scheduler model.PaymentScheduler) (model.PaymentScheduler, error)
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID)
	FindById(id uuid.UUID) (model.PaymentScheduler, error)
	FindByHouseId(houseId uuid.UUID) []model.PaymentScheduler
	FindByUserId(userId uuid.UUID) []model.PaymentScheduler
}

func (p *PaymentSchedulerRepositoryObject) Create(scheduler model.PaymentScheduler) (model.PaymentScheduler, error) {
	return scheduler, p.database.Create(&scheduler)
}

func (p *PaymentSchedulerRepositoryObject) ExistsById(id uuid.UUID) bool {
	return p.database.ExistsById(model.PaymentScheduler{}, id)
}

func (p *PaymentSchedulerRepositoryObject) DeleteById(id uuid.UUID) {
	p.database.D().Delete(model.PaymentScheduler{}, id)
}

func (p *PaymentSchedulerRepositoryObject) FindById(id uuid.UUID) (response model.PaymentScheduler, err error) {
	return response, p.database.FindById(&response, id)
}

func (p *PaymentSchedulerRepositoryObject) FindByHouseId(houseId uuid.UUID) (response []model.PaymentScheduler) {
	if tx := p.database.D().Find(&response, "house_id = ?", houseId); tx.Error != nil {
		log.Error().Err(tx.Error).Msg("")
	}
	return response
}

func (p *PaymentSchedulerRepositoryObject) FindByUserId(userId uuid.UUID) (response []model.PaymentScheduler) {
	if tx := p.database.D().Find(&response, "user_id = ?", userId); tx.Error != nil {
		log.Error().Err(tx.Error).Msg("")
	}
	return response
}
