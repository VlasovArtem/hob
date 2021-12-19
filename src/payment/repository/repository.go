package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/google/uuid"
	"reflect"
)

var PaymentRepositoryType = reflect.TypeOf(PaymentRepositoryObject{})
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

func (p *PaymentRepositoryObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewPaymentRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (p *PaymentRepositoryObject) GetEntity() interface{} {
	return entity
}

type PaymentRepository interface {
	Create(entity model.Payment) (model.Payment, error)
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) (model.Payment, error)
	FindByHouseId(houseId uuid.UUID) []model.PaymentDto
	FindByUserId(userId uuid.UUID) []model.PaymentDto
	FindByProviderId(providerId uuid.UUID) []model.PaymentDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(entity model.Payment) error
}

func (p *PaymentRepositoryObject) Create(entity model.Payment) (model.Payment, error) {
	return entity, p.database.Create(&entity)
}

func (p *PaymentRepositoryObject) Delete(id uuid.UUID) (err error) {
	return p.database.Delete(id)
}

func (p *PaymentRepositoryObject) FindById(id uuid.UUID) (response model.Payment, err error) {
	return response, p.database.Find(&response, id)
}

func (p *PaymentRepositoryObject) FindByHouseId(houseId uuid.UUID) (response []model.PaymentDto) {
	return p.findBy("house_id = ?", houseId)
}

func (p *PaymentRepositoryObject) FindByUserId(userId uuid.UUID) (response []model.PaymentDto) {
	return p.findBy("user_id = ?", userId)
}

func (p *PaymentRepositoryObject) FindByProviderId(providerId uuid.UUID) []model.PaymentDto {
	return p.findBy("provider_id = ?", providerId)
}

func (p *PaymentRepositoryObject) findBy(query interface{}, conditions ...interface{}) (response []model.PaymentDto) {
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
