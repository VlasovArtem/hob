package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/db"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
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
		AddAfterTest(truncateDynamic).
		AddAfterSuite(func(service db.DatabaseService) {
			database.TruncateTable(service, houseModel.House{})
			database.TruncateTable(service, userModel.User{})
		}).
		ExecuteMigration(userModel.User{}, groupModel.Group{}, houseModel.House{}, model.Income{})

	i.createdUser = userMocks.GenerateUser()
	i.CreateEntity(&i.createdUser)

	i.createdHouse = houseMocks.GenerateHouse(i.createdUser.Id)
	i.CreateEntity(&i.createdHouse)
}

func TestIncomeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeRepositoryTestSuite))
}

func (i *IncomeRepositoryTestSuite) Test_Create() {
	income := mocks.GenerateIncome(&i.createdHouse.Id)

	actual, err := i.repository.Create(income)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), income, actual)
}

func (i *IncomeRepositoryTestSuite) Test_Create_WithMissingHouse() {
	houseId := uuid.New()
	income := mocks.GenerateIncome(&houseId)

	actual, err := i.repository.Create(income)

	assert.NotNil(i.T(), err)
	assert.Equal(i.T(), income, actual)
}

func (i *IncomeRepositoryTestSuite) Test_Create_WithGroupId() {
	group := i.createGroup()

	income := mocks.GenerateIncome(nil)
	income.Groups = []groupModel.Group{group}

	actual, err := i.repository.Create(income)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), income, actual)
}

func (i *IncomeRepositoryTestSuite) Test_CreateBatch() {
	first := mocks.GenerateIncome(&i.createdHouse.Id)
	first.Name = "First Income"
	second := mocks.GenerateIncome(&i.createdHouse.Id)
	second.Name = "Second Income"

	actual, err := i.repository.CreateBatch([]model.Income{first, second})

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.Income{first, second}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindById() {
	income := i.createIncome()

	actual, err := i.repository.FindById(income.Id)

	income.Groups = []groupModel.Group{}
	assert.Nil(i.T(), err)
	assert.Equal(i.T(), income, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindById_WithMissingId() {
	actual, err := i.repository.FindById(uuid.New())

	assert.ErrorIs(i.T(), err, gorm.ErrRecordNotFound)
	assert.Equal(i.T(), model.Income{}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindByHouseId() {
	income := i.createIncomeWithHouse()

	actual, err := i.repository.FindByHouseId(*income.HouseId)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{income.ToDto()}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindByHouseId_WithMissingId() {
	actual, err := i.repository.FindByHouseId(uuid.New())

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindByGroupIds() {
	first := i.createGroup()
	second := i.createGroup()
	third := i.createGroup()
	firstIncome := i.createIncomeWithGroups([]groupModel.Group{first})
	secondIncome := i.createIncomeWithGroups([]groupModel.Group{second})
	i.createIncomeWithGroups([]groupModel.Group{third})

	actual, err := i.repository.FindByGroupIds([]uuid.UUID{first.Id, second.Id})

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{
		firstIncome.ToDto(),
		secondIncome.ToDto(),
	}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindByGroupIds_WithNotMatchingIds() {
	i.createIncome()

	actual, err := i.repository.FindByGroupIds([]uuid.UUID{uuid.New()})

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_ExistsById() {
	income := i.createIncome()

	assert.True(i.T(), i.repository.ExistsById(income.Id))
}

func (i *IncomeRepositoryTestSuite) Test_ExistsById_WithMissingId() {
	assert.False(i.T(), i.repository.ExistsById(uuid.New()))
}

func (i *IncomeRepositoryTestSuite) Test_DeleteById() {
	income := i.createIncome()

	assert.Nil(i.T(), i.repository.DeleteById(income.Id))
}

func (i *IncomeRepositoryTestSuite) Test_DeleteById_WithMissingId() {
	assert.Nil(i.T(), i.repository.DeleteById(uuid.New()))
}

func (i *IncomeRepositoryTestSuite) Test_Update() {
	income := i.createIncome()

	updatedIncome := model.UpdateIncomeRequest{
		Name:        fmt.Sprintf("%s-new", income.Name),
		Description: fmt.Sprintf("%s-new", income.Description),
		Date:        mocks.Date,
		Sum:         income.Sum + 100.0,
	}

	err := i.repository.Update(income.Id, updatedIncome)

	assert.Nil(i.T(), err)

	response, err := i.repository.FindById(income.Id)
	assert.Nil(i.T(), err)
	assert.Equal(i.T(), model.Income{
		Id:          income.Id,
		Name:        "Name-new",
		Description: "Description-new",
		Date:        updatedIncome.Date,
		Sum:         200.1,
		HouseId:     income.HouseId,
		House:       income.House,
		Groups:      []groupModel.Group{},
	}, response)
}

func (i *IncomeRepositoryTestSuite) createIncomeWithHouse() model.Income {
	createdHouse := houseMocks.GenerateHouse(i.createdUser.Id)
	i.CreateEntity(createdHouse)

	income := mocks.GenerateIncome(&createdHouse.Id)
	i.CreateEntity(income)

	return income
}

func (i *IncomeRepositoryTestSuite) createIncome() model.Income {
	income := mocks.GenerateIncome(&i.createdHouse.Id)

	i.CreateEntity(income)

	return income
}

func (i *IncomeRepositoryTestSuite) createIncomeWithGroups(groups []groupModel.Group) (income model.Income) {
	income = mocks.GenerateIncome(&i.createdHouse.Id)
	income.Groups = groups

	create, err := i.repository.Create(income)

	assert.Nil(i.T(), err)

	return create
}

func truncateDynamic(service db.DatabaseService) {
	service.D().Exec("DELETE FROM income_groups")
	database.TruncateTable(service, groupModel.Group{})
	database.TruncateTable(service, model.Income{})
}

func (i *IncomeRepositoryTestSuite) createGroup() groupModel.Group {
	id := uuid.New()
	group := groupModel.Group{
		Id:      id,
		Name:    "Test Group " + id.String(),
		OwnerId: i.createdUser.Id,
	}

	i.CreateEntity(group)

	return group
}
