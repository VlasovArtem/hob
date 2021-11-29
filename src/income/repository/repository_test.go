package repository

import (
	"github.com/VlasovArtem/hob/src/db"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/income/mocks"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type IncomeRepositoryTestSuite struct {
	database.DBTestSuite
	repository   IncomeRepository
	createdUser  userModel.User
	createdHouse houseModel.House
}

func (i *IncomeRepositoryTestSuite) SetupSuite() {
	i.InitDBTestSuite()

	i.CreateRepository(
		func(service db.DatabaseService) {
			i.repository = NewIncomeRepository(service)
		},
	).
		AddMigrators(userModel.User{}, houseModel.House{}, model.Income{})

	i.createdUser = userMocks.GenerateUser()
	i.CreateEntity(&i.createdUser)

	i.createdHouse = houseMocks.GenerateHouse(i.createdUser.Id)
	i.CreateEntity(&i.createdHouse)
}

func (i *IncomeRepositoryTestSuite) TearDownSuite() {
	i.TearDown()
}

func TestIncomeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeRepositoryTestSuite))
}

func (i *IncomeRepositoryTestSuite) Test_Create() {
	income := mocks.GenerateIncome(i.createdHouse.Id)

	actual, err := i.repository.Create(income)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), income, actual)
}

func (i *IncomeRepositoryTestSuite) Test_Creat_WithMissingHouse() {
	income := mocks.GenerateIncome(uuid.New())

	actual, err := i.repository.Create(income)

	assert.NotNil(i.T(), err)
	assert.Equal(i.T(), income, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindResponseById() {
	income := i.createIncome()

	actual, err := i.repository.FindResponseById(income.Id)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), income.ToResponse(), actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindResponseById_WithMissingId() {
	actual, err := i.repository.FindResponseById(uuid.New())

	assert.ErrorIs(i.T(), err, gorm.ErrRecordNotFound)
	assert.Equal(i.T(), model.IncomeResponse{}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindResponseByHouseId() {
	income := i.createIncome()

	actual := i.repository.FindResponseByHouseId(income.HouseId)

	var actualResponse model.IncomeResponse

	for _, response := range actual {
		if response.Id == income.Id {
			actualResponse = response
			break
		}
	}
	assert.Equal(i.T(), income.ToResponse(), actualResponse)
}

func (i *IncomeRepositoryTestSuite) Test_FindResponseByHouseId_WithMissingId() {
	actual := i.repository.FindResponseByHouseId(uuid.New())

	assert.Equal(i.T(), []model.IncomeResponse{}, actual)
}

func (i *IncomeRepositoryTestSuite) createIncome() model.Income {
	income := mocks.GenerateIncome(i.createdHouse.Id)

	create, err := i.repository.Create(income)

	assert.Nil(i.T(), err)

	return create
}
