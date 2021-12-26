package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/google/uuid"
	"reflect"
)

var (
	MeterRepositoryType = reflect.TypeOf(MeterRepositoryObject{})
	entity              = model.Meter{}
)

type MeterRepositoryObject struct {
	database db.ModeledDatabase
}

func NewMeterRepository(database db.DatabaseService) MeterRepository {
	return &MeterRepositoryObject{
		database: db.ModeledDatabase{
			DatabaseService: database,
			Model:           entity,
		},
	}
}

func (m *MeterRepositoryObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewMeterRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (m *MeterRepositoryObject) GetEntity() interface{} {
	return entity
}

type MeterRepository interface {
	Create(entity model.Meter) (model.Meter, error)
	ExistsById(id uuid.UUID) bool
	FindById(id uuid.UUID) (model.Meter, error)
	FindByPaymentId(paymentId uuid.UUID) (model.Meter, error)
	FindByHouseId(houseId uuid.UUID) []model.Meter
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, meter model.Meter) error
}

func (m *MeterRepositoryObject) Create(entity model.Meter) (model.Meter, error) {
	if err := m.database.Create(&entity); err != nil {
		return entity, err
	}
	return entity, nil
}

func (m *MeterRepositoryObject) ExistsById(id uuid.UUID) bool {
	return m.database.Exists(id)
}

func (m *MeterRepositoryObject) FindById(id uuid.UUID) (response model.Meter, err error) {
	return response, m.database.Find(&response, id)
}

func (m *MeterRepositoryObject) FindByPaymentId(id uuid.UUID) (response model.Meter, err error) {
	return response, m.database.FirstBy(&response, "payment_id = ?", id)
}

func (m *MeterRepositoryObject) FindByHouseId(houseId uuid.UUID) (response []model.Meter) {
	_ = m.database.FindBy(&response, "house_id = ?", houseId)

	return response
}

func (m *MeterRepositoryObject) DeleteById(id uuid.UUID) error {
	return m.database.Delete(id)
}

func (m *MeterRepositoryObject) Update(id uuid.UUID, entity model.Meter) error {
	return m.database.Update(id, entity, "HouseId", "House", "PaymentId", "Payment")
}
