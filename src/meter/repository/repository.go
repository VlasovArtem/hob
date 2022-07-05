package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MeterRepositoryStr struct {
	db.ModeledDatabase[model.Meter]
}

func NewMeterRepository(database db.DatabaseService) MeterRepository {
	return &MeterRepositoryStr{db.NewModeledDatabase(model.Meter{}, database)}
}

func (m *MeterRepositoryStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewMeterRepository(factory.FindRequiredByObject(db.Database{}).(db.DatabaseService))
}

type MeterRepository interface {
	transactional.Transactional[MeterRepository]
	db.ModeledDatabase[model.Meter]
	FindByPaymentId(paymentId uuid.UUID) (model.Meter, error)
}

func (m *MeterRepositoryStr) FindByPaymentId(id uuid.UUID) (response model.Meter, err error) {
	return m.FirstBy("payment_id = ?", id)
}

func (m *MeterRepositoryStr) Transactional(tx *gorm.DB) MeterRepository {
	return &MeterRepositoryStr{
		ModeledDatabase: db.NewTransactionalModeledDatabase(m.GetEntity(), tx),
	}
}
