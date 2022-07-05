package repository

import (
	"github.com/VlasovArtem/hob/src/db"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	incomeModel "github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/scheduler/mocks"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	"github.com/VlasovArtem/hob/src/scheduler"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type IncomeSchedulerRepositoryTestSuite struct {
	database.DBTestSuite[model.IncomeScheduler]
	repository   IncomeSchedulerRepository
	createdUser  userModel.User
	createdHouse houseModel.House
}

func (i *IncomeSchedulerRepositoryTestSuite) SetupSuite() {
	i.InitDBTestSuite()

	i.CreateRepository(
		func(service db.ModeledDatabase[model.IncomeScheduler]) {
			i.repository = NewIncomeSchedulerRepository(service)
		},
	).
		AddAfterTest(func(service db.ModeledDatabase[model.IncomeScheduler]) {
			testhelper.TruncateTable(service, model.IncomeScheduler{})
		}).
		AddAfterSuite(func(service db.ModeledDatabase[model.IncomeScheduler]) {
			testhelper.TruncateTable(service, houseModel.House{})
			testhelper.TruncateTable(service, userModel.User{})
		}).
		ExecuteMigration(userModel.User{}, houseModel.House{}, model.IncomeScheduler{})

	i.createdUser = userMocks.GenerateUser()
	i.CreateEntity(&i.createdUser)

	i.createdHouse = houseMocks.GenerateHouse(i.createdUser.Id)
	i.CreateEntity(&i.createdHouse)
}

func TestPaymentRepositorySchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeSchedulerRepositoryTestSuite))
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_FindByHouseId() {
	payment := i.createIncomeSchedulerWithNewHouse()

	actual, err := i.repository.FindByHouseId(*payment.HouseId)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeSchedulerDto{payment.ToDto()}, actual)
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_FindByHouseId_WithMissingId() {
	actual, err := i.repository.FindByHouseId(uuid.New())

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeSchedulerDto{}, actual)
}

func (i *IncomeSchedulerRepositoryTestSuite) Test_Update() {
	incomeScheduler := i.createIncomeScheduler()

	updatedIncomeScheduler := model.UpdateIncomeSchedulerRequest{
		Name:        "New Name",
		Description: "New Description",
		Sum:         1010,
		Spec:        scheduler.WEEKLY,
	}

	err := i.repository.Update(incomeScheduler.Id, updatedIncomeScheduler)

	assert.Nil(i.T(), err)

	actual, err := i.repository.Find(incomeScheduler.Id)

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
