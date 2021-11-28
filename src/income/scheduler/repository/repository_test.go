package repository

import (
	"db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	houseMocks "house/mocks"
	houseModel "house/model"
	"income/scheduler/mocks"
	"income/scheduler/model"
	"log"
	"test/testhelper/database"
	"testing"
	userMocks "user/mocks"
	userModel "user/model"
)

type IncomeSchedulerRepositoryTestSuite struct {
	suite.Suite
	database     db.DatabaseService
	repository   IncomeSchedulerRepository
	createdUser  userModel.User
	createdHouse houseModel.House
}

func (i *IncomeSchedulerRepositoryTestSuite) SetupSuite() {
	config := db.NewDefaultDatabaseConfiguration()
	config.DBName = "hob_test"
	i.database = db.NewDatabaseService(config)
	i.repository = NewIncomeSchedulerRepository(i.database)
	err := i.database.D().AutoMigrate(model.IncomeScheduler{})

	if err != nil {
		log.Fatal(err)
	}

	i.createdUser = userMocks.GenerateUser()
	err = i.database.Create(&i.createdUser)

	if err != nil {
		log.Fatal(err)
	}
}

func (i *IncomeSchedulerRepositoryTestSuite) TearDownSuite() {
	database.DropTable(i.database.D(), houseModel.House{})
	database.DropTable(i.database.D(), userModel.User{})
	database.DropTable(i.database.D(), model.IncomeScheduler{})
}

func (i *IncomeSchedulerRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	i.createdHouse = houseMocks.GenerateHouse(i.createdUser.Id)

	err := i.database.Create(&i.createdHouse)

	if err != nil {
		log.Fatal(err)
	}
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
