package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/google/uuid"
	"reflect"
)

var PaymentSchedulerRepositoryType = reflect.TypeOf(PaymentSchedulerRepositoryObject{})
var entity = model.PaymentScheduler{}

type PaymentSchedulerRepositoryObject struct {
	database db.ModeledDatabase
}

func NewPaymentSchedulerRepository(database db.DatabaseService) PaymentSchedulerRepository {
	return &PaymentSchedulerRepositoryObject{
		db.ModeledDatabase{
			DatabaseService: database,
			Model:           entity,
		},
	}
}

func (p *PaymentSchedulerRepositoryObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewPaymentSchedulerRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (p *PaymentSchedulerRepositoryObject) GetEntity() any {
	return entity
}

type PaymentSchedulerRepository interface {
	Create(scheduler model.PaymentScheduler) (model.PaymentScheduler, error)
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID)
	FindById(id uuid.UUID) (model.PaymentScheduler, error)
	FindByHouseId(houseId uuid.UUID) []model.PaymentSchedulerDto
	FindByUserId(userId uuid.UUID) []model.PaymentSchedulerDto
	FindByProviderId(providerId uuid.UUID) []model.PaymentSchedulerDto
	Update(id uuid.UUID, request model.UpdatePaymentSchedulerRequest) (model.PaymentScheduler, error)
}

func (p *PaymentSchedulerRepositoryObject) Create(scheduler model.PaymentScheduler) (model.PaymentScheduler, error) {
	return scheduler, p.database.Create(&scheduler)
}

func (p *PaymentSchedulerRepositoryObject) ExistsById(id uuid.UUID) bool {
	return p.database.Exists(id)
}

func (p *PaymentSchedulerRepositoryObject) DeleteById(id uuid.UUID) {
	_ = p.database.Delete(id)
}

func (p *PaymentSchedulerRepositoryObject) FindById(id uuid.UUID) (response model.PaymentScheduler, err error) {
	return response, p.database.FindById(&response, id)
}

func (p *PaymentSchedulerRepositoryObject) FindByHouseId(houseId uuid.UUID) (response []model.PaymentSchedulerDto) {
	return p.findBy("house_id = ?", houseId)
}

func (p *PaymentSchedulerRepositoryObject) FindByUserId(userId uuid.UUID) (response []model.PaymentSchedulerDto) {
	return p.findBy("user_id = ?", userId)
}

func (p *PaymentSchedulerRepositoryObject) FindByProviderId(providerId uuid.UUID) (response []model.PaymentSchedulerDto) {
	return p.findBy("provider_id = ?", providerId)
}

func (p *PaymentSchedulerRepositoryObject) findBy(query any, conditions ...any) (response []model.PaymentSchedulerDto) {
	_ = p.database.FindBy(&response, query, conditions...)
	return response
}

func (p *PaymentSchedulerRepositoryObject) Update(id uuid.UUID, request model.UpdatePaymentSchedulerRequest) (response model.PaymentScheduler, error error) {
	if err := p.database.Update(id, request); err != nil {
		return response, err
	}

	return p.FindById(id)
}
