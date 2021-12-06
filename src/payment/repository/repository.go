package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/google/uuid"
)

type PaymentRepositoryObject struct {
	database db.DatabaseService
}

func NewPaymentRepository(service db.DatabaseService) PaymentRepository {
	return &PaymentRepositoryObject{service}
}

func (p *PaymentRepositoryObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return NewPaymentRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (p *PaymentRepositoryObject) GetEntity() interface{} {
	return model.Payment{}
}

type PaymentRepository interface {
	Create(entity model.Payment) (model.Payment, error)
	FindById(id uuid.UUID) (model.Payment, error)
	FindByHouseId(houseId uuid.UUID) []model.Payment
	FindByUserId(userId uuid.UUID) []model.Payment
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(entity model.Payment) error
}

func (p *PaymentRepositoryObject) Create(entity model.Payment) (model.Payment, error) {
	return entity, p.database.Create(&entity)
}

func (p *PaymentRepositoryObject) FindById(id uuid.UUID) (response model.Payment, err error) {
	return response, p.database.FindById(&response, id)
}

func (p *PaymentRepositoryObject) FindByHouseId(houseId uuid.UUID) (response []model.Payment) {
	p.database.DM(model.Payment{}).Where("house_id = ?", houseId).Find(&response)

	return response
}

func (p *PaymentRepositoryObject) FindByUserId(userId uuid.UUID) (response []model.Payment) {
	p.database.DM(model.Payment{}).Where("user_id = ?", userId).Find(&response)

	return response
}

func (p *PaymentRepositoryObject) ExistsById(id uuid.UUID) bool {
	return p.database.ExistsById(model.Payment{}, id)
}

func (p *PaymentRepositoryObject) DeleteById(id uuid.UUID) error {
	return p.database.D().Delete(model.Payment{}, id).Error
}

func (p *PaymentRepositoryObject) Update(entity model.Payment) error {
	return p.database.UpdateById(model.Payment{}, entity.Id, entity, "HouseId", "House", "UserId", "User")
}
