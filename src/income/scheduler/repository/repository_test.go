package repository

import (
	"github.com/VlasovArtem/hob/src/db"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/income/scheduler/mocks"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
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
	i.CreateEntity(&i.createdUser)

	i.AddBeforeTest(
		func(service db.DatabaseService) {
			i.createdHouse = houseMocks.GenerateHouse(i.createdUser.Id)
			i.CreateEntity(&i.createdHouse)
		})
}

func (i *IncomeSchedulerRepositoryTestSuite) TearDownSuite() {
	i.TearDown()
}

func TestPaymentRepositorySchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeSchedulerRepositoryTestSuite))
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_Create() {
	incomeScheduler := mocks.GenerateIncomeScheduler(i.createdHouse.Id)

	actual, err := i.repository.Create(incomeScheduler)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), incomeScheduler, actual)
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
	payment := i.createIncomeScheduler()

	actual := i.repository.FindByHouseId(payment.HouseId)

	assert.Equal(i.T(), []model.IncomeScheduler{payment}, actual)
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_FindByHouseId_WithMissingId() {
	actual := i.repository.FindByHouseId(uuid.New())

	assert.Equal(i.T(), []model.IncomeScheduler{}, actual)
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

	i.repository.DeleteById(payment.Id)

	assert.False(i.T(), i.repository.ExistsById(payment.Id))
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_DeleteById_WithMissingId() {
	i.repository.DeleteById(uuid.New())
}

func (i *IncomeSchedulerRepositoryTestSuite) createIncomeScheduler() model.IncomeScheduler {
	payment := mocks.GenerateIncomeScheduler(i.createdHouse.Id)

	create, err := i.repository.Create(payment)

	assert.Nil(i.T(), err)

	return create
}
