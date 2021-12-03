package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/google/uuid"
)

type MeterRepositoryObject struct {
	database db.DatabaseService
}

func NewMeterRepository(database db.DatabaseService) MeterRepository {
	return &MeterRepositoryObject{database}
}

func (m *MeterRepositoryObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return NewMeterRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (m *MeterRepositoryObject) GetEntity() interface{} {
	return model.Meter{}
}

type MeterRepository interface {
	Create(entity model.Meter) (model.Meter, error)
	ExistsById(id uuid.UUID) bool
	FindById(id uuid.UUID) (model.Meter, error)
	FindByPaymentId(paymentId uuid.UUID) (model.Meter, error)
	FindByHouseId(houseId uuid.UUID) []model.Meter
}

func (m *MeterRepositoryObject) Create(entity model.Meter) (model.Meter, error) {
	if err := m.database.Create(&entity); err != nil {
		return entity, err
	}
	return entity, nil
}

func (m *MeterRepositoryObject) ExistsById(id uuid.UUID) bool {
	return m.database.ExistsById(model.Meter{}, id)
}

func (m *MeterRepositoryObject) FindById(id uuid.UUID) (response model.Meter, err error) {
	return response, m.database.FindById(&response, id)
}

func (m *MeterRepositoryObject) FindByPaymentId(id uuid.UUID) (response model.Meter, err error) {
	if tx := m.database.DM(model.Meter{}).First(&response, "payment_id = ?", id); tx.Error != nil {
		return response, tx.Error
	}
	return response, err
}

func (m *MeterRepositoryObject) FindByHouseId(houseId uuid.UUID) (response []model.Meter) {
	m.database.DM(model.Meter{}).Where("house_id = ?", houseId).Find(&response)

	return response
}
