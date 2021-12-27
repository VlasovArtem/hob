package repository

import (
	"github.com/VlasovArtem/hob/src/db"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	incomeModel "github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/scheduler/mocks"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	"github.com/VlasovArtem/hob/src/scheduler"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type IncomeSchedulerRepositoryTestSuite struct {
	database.DBTestSuite
	repository   IncomeSchedulerRepository
	createdUser  userModel.User
	createdHouse houseModel.House
}

func (i *IncomeSchedulerRepositoryTestSuite) SetupSuite() {
	i.InitDBTestSuite()

	i.CreateRepository(
		func(service db.DatabaseService) {
			i.repository = NewIncomeSchedulerRepository(service)
		},
	).
		AddMigrators(userModel.User{}, houseModel.House{}, model.IncomeScheduler{})

	i.createdUser = userMocks.GenerateUser()
	i.CreateConstantEntity(&i.createdUser)

	i.createdHouse = houseMocks.GenerateHouse(i.createdUser.Id)
	i.CreateConstantEntity(&i.createdHouse)
}

func TestPaymentRepositorySchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeSchedulerRepositoryTestSuite))
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_Create() {
	incomeScheduler := mocks.GenerateIncomeScheduler(i.createdHouse.Id)

	actual, err := i.repository.Create(incomeScheduler)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), incomeScheduler, actual)

	i.Delete(incomeScheduler)
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_Creat_WithMissingHouse() {
	incomeScheduler := mocks.GenerateIncomeScheduler(uuid.New())

	actual, err := i.repository.Create(incomeScheduler)

	assert.NotNil(i.T(), err)
	assert.Equal(i.T(), incomeScheduler, actual)
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_FindById() {
	payment := i.createIncomeScheduler()

	actual, err := i.repository.FindById(payment.Id)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), payment, actual)
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_FindById_WithMissingId() {
	actual, err := i.repository.FindById(uuid.New())

	assert.ErrorIs(i.T(), err, gorm.ErrRecordNotFound)
	assert.Equal(i.T(), model.IncomeScheduler{}, actual)
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_FindByHouseId() {
	payment := i.createIncomeSchedulerWithNewHouse()

	actual, err := i.repository.FindByHouseId(payment.HouseId)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeSchedulerDto{payment.ToDto()}, actual)
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_FindByHouseId_WithMissingId() {
	actual, err := i.repository.FindByHouseId(uuid.New())

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeSchedulerDto{}, actual)
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_ExistsById() {
	payment := i.createIncomeScheduler()

	assert.True(i.T(), i.repository.ExistsById(payment.Id))
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_ExistsById_WithMissingId() {
	assert.False(i.T(), i.repository.ExistsById(uuid.New()))
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_DeleteById() {
	payment := i.createIncomeScheduler()

	assert.True(i.T(), i.repository.ExistsById(payment.Id))

	err := i.repository.DeleteById(payment.Id)

	assert.Nil(i.T(), err)

	assert.False(i.T(), i.repository.ExistsById(payment.Id))
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_DeleteById_WithMissingId() {
	err := i.repository.DeleteById(uuid.New())

	assert.Nil(i.T(), err)
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_Update() {
	incomeScheduler := i.createIncomeScheduler()

	updatedIncomeScheduler := model.UpdateIncomeSchedulerRequest{
		Name:        "New Name",
		Description: "New Description",
		Sum:         1010,
		Spec:        scheduler.WEEKLY,
	}

	actual, err := i.repository.Update(incomeScheduler.Id, updatedIncomeScheduler)

	assert.Nil(i.T(), err)

	assert.Equal(i.T(), model.IncomeScheduler{
		Income: incomeModel.Income{
			Id:          incomeScheduler.Id,
			Name:        "New Name",
			Description: "New Description",
			Date:        incomeScheduler.Date,
			Sum:         1010,
			HouseId:     incomeScheduler.HouseId,
			House:       incomeScheduler.House,
		},
		Spec: scheduler.WEEKLY,
	}, actual)
}

func (i *IncomeSchedulerRepositoryTestSuite) createIncomeSchedulerWithNewHouse() model.IncomeScheduler {
	createdHouse := houseMocks.GenerateHouse(i.createdUser.Id)
	i.CreateEntity(&createdHouse)

	payment := mocks.GenerateIncomeScheduler(createdHouse.Id)

	i.CreateEntity(payment)

	return payment
}

func (i *IncomeSchedulerRepositoryTestSuite) createIncomeScheduler() model.IncomeScheduler {
	payment := mocks.GenerateIncomeScheduler(i.createdHouse.Id)

	i.CreateEntity(payment)

	return payment
}
