package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/db"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/income/mocks"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type IncomeRepositoryTestSuite struct {
	database.DBTestSuite[model.Income]
	repository   IncomeRepository
	createdUser  userModel.User
	createdHouse houseModel.House
}

func (i *IncomeRepositoryTestSuite) SetupSuite() {
	i.InitDBTestSuite()

	i.CreateRepository(
		func(service db.ModeledDatabase[model.Income]) {
			i.repository = NewIncomeRepository(service)
		},
	).
		AddAfterTest(func(service db.ModeledDatabase[model.Income]) {
			service.DB().Exec("DELETE FROM income_groups")
			service.DB().Exec("DELETE FROM house_groups")
			testhelper.TruncateTable(service, groupModel.Group{})
			testhelper.TruncateTable(service, model.Income{})
		}).
		AddAfterSuite(func(service db.ModeledDatabase[model.Income]) {
			service.DB().Exec("DELETE FROM income_groups")
			service.DB().Exec("DELETE FROM house_groups")
			testhelper.TruncateTable(service, groupModel.Group{})
			testhelper.TruncateTable(service, model.Income{})
			testhelper.TruncateTable(service, houseModel.House{})
			testhelper.TruncateTable(service, userModel.User{})
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

func (i *IncomeRepositoryTestSuite) Test_CreateBatch() {
	first := mocks.GenerateIncome(&i.createdHouse.Id)
	first.Name = "First Income"
	second := mocks.GenerateIncome(&i.createdHouse.Id)
	second.Name = "Second Income"

	incomes := []model.Income{first, second}
	err := i.repository.Create(&incomes)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.Income{first, second}, incomes)
}

func (i *IncomeRepositoryTestSuite) Test_FindByHouseId() {
	income := i.createIncomeWithHouse()

	actual, err := i.repository.FindByHouseId(*income.HouseId, 10, 0, nil, nil)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{income.ToDto()}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindByHouseId_WithFromAndTo() {
	first := mocks.GenerateIncome(&i.createdHouse.Id)
	first.Date = time.Now().AddDate(0, 0, -1)
	i.CreateEntity(&first)

	second := mocks.GenerateIncome(&i.createdHouse.Id)
	second.Date = time.Now().Truncate(time.Microsecond)
	i.CreateEntity(&second)

	from := time.Now().Add(-time.Hour * 12)
	to := time.Now()

	actual, err := i.repository.FindByHouseId(i.createdHouse.Id, 10, 0, &from, &to)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{second.ToDto()}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindByHouseId_WithFrom() {
	first := mocks.GenerateIncome(&i.createdHouse.Id)
	first.Date = time.Now().AddDate(0, 0, -1)
	i.CreateEntity(&first)

	second := mocks.GenerateIncome(&i.createdHouse.Id)
	second.Date = time.Now().Truncate(time.Microsecond)
	i.CreateEntity(&second)

	from := time.Now().Add(-time.Hour * 12)

	actual, err := i.repository.FindByHouseId(i.createdHouse.Id, 10, 0, &from, nil)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{second.ToDto()}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindByHouseId_WithGroupsAndHouseId() {
	group := i.createGroup()
	house := houseMocks.GenerateHouse(i.createdUser.Id)
	house.Groups = []groupModel.Group{group}
	i.CreateEntity(&house)

	incomeWithHouseId := mocks.GenerateIncome(&house.Id)
	incomeWithHouseId.Date = time.Now().Truncate(time.Microsecond)
	i.CreateEntity(&incomeWithHouseId)

	incomeWithGroups := mocks.GenerateIncome(nil)
	incomeWithGroups.Groups = []groupModel.Group{group}
	incomeWithGroups.Date = time.Now().Truncate(time.Microsecond)
	i.CreateEntity(&incomeWithGroups)

	actual, err := i.repository.FindByHouseId(house.Id, 10, 0, nil, nil)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{incomeWithGroups.ToDto(), incomeWithHouseId.ToDto()}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindByHouseId_WithMissingId() {
	actual, err := i.repository.FindByHouseId(uuid.New(), 10, 0, nil, nil)

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

	actual, err := i.repository.FindByGroupIds([]uuid.UUID{first.Id, second.Id}, 10, 0, nil, nil)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{
		secondIncome.ToDto(),
		firstIncome.ToDto(),
	}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindByGroupIds_WithFromAndTo() {
	first := i.createGroup()
	firstIncome := mocks.GenerateIncome(&i.createdHouse.Id)
	firstIncome.Date = time.Now().AddDate(0, 0, -1).Truncate(time.Microsecond)
	firstIncome.Groups = []groupModel.Group{first}
	i.CreateEntity(&firstIncome)

	second := i.createGroup()
	secondIncome := mocks.GenerateIncome(&i.createdHouse.Id)
	secondIncome.Date = time.Now().Truncate(time.Microsecond)
	secondIncome.Groups = []groupModel.Group{second}
	i.CreateEntity(&secondIncome)

	i.createIncomeWithGroups([]groupModel.Group{i.createGroup()})

	from := time.Now().Add(-time.Hour * 12)
	to := time.Now()

	actual, err := i.repository.FindByGroupIds([]uuid.UUID{first.Id, second.Id}, 10, 0, &from, &to)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{secondIncome.ToDto()}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindByGroupIds_WithFrom() {
	first := i.createGroup()
	firstIncome := mocks.GenerateIncome(&i.createdHouse.Id)
	firstIncome.Date = time.Now().AddDate(0, 0, -1)
	firstIncome.Groups = []groupModel.Group{first}
	i.CreateEntity(&firstIncome)

	second := i.createGroup()
	secondIncome := mocks.GenerateIncome(&i.createdHouse.Id)
	secondIncome.Date = time.Now().Truncate(time.Microsecond)
	secondIncome.Groups = []groupModel.Group{second}
	i.CreateEntity(&secondIncome)

	i.createIncomeWithGroups([]groupModel.Group{i.createGroup()})

	from := time.Now().Add(-time.Hour * 12)

	actual, err := i.repository.FindByGroupIds([]uuid.UUID{first.Id, second.Id}, 10, 0, &from, nil)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{secondIncome.ToDto()}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_FindByGroupIds_WithNotMatchingIds() {
	i.createIncome()

	actual, err := i.repository.FindByGroupIds([]uuid.UUID{uuid.New()}, 10, 0, nil, nil)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{}, actual)
}

func (i *IncomeRepositoryTestSuite) Test_Update() {
	income := i.createIncome()

	updatedIncome := model.UpdateIncomeRequest{
		Name:        fmt.Sprintf("%s-new", income.Name),
		Description: fmt.Sprintf("%s-new", income.Description),
		Date:        mocks.Date,
		Sum:         income.Sum + 100.0,
	}

	err := i.repository.UpdateByRequest(income.Id, updatedIncome)

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

func (i *IncomeRepositoryTestSuite) Test_UpdateByRequest_WithAddingGroups() {
	group := groupModel.Group{
		Id:      uuid.New(),
		Name:    "Test Group",
		OwnerId: i.createdUser.Id,
	}
	i.CreateEntity(group)

	income := i.createIncome()

	updatedHouse := model.UpdateIncomeRequest{
		GroupIds: []uuid.UUID{group.Id},
	}

	err := i.repository.UpdateByRequest(income.Id, updatedHouse)

	assert.Nil(i.T(), err)

	response, err := i.repository.FindById(income.Id)
	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []groupModel.Group{group}, response.Groups)
}

func (i *IncomeRepositoryTestSuite) Test_UpdateByRequest_WithUpdatingGroups() {
	group := groupModel.Group{
		Id:      uuid.New(),
		Name:    "Test Group",
		OwnerId: i.createdUser.Id,
	}

	income := i.createIncomeWithGroups([]groupModel.Group{group})

	newGroup := groupModel.Group{
		Id:      uuid.New(),
		Name:    "New Test Group",
		OwnerId: i.createdUser.Id,
	}

	i.CreateEntity(newGroup)

	updatedHouse := model.UpdateIncomeRequest{
		GroupIds: []uuid.UUID{newGroup.Id},
	}

	err := i.repository.UpdateByRequest(income.Id, updatedHouse)

	assert.Nil(i.T(), err)

	response, err := i.repository.FindById(income.Id)
	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []groupModel.Group{newGroup}, response.Groups)
}

func (i *IncomeRepositoryTestSuite) Test_UpdateByRequest_WithExtendingGroups() {
	group := groupModel.Group{
		Id:      uuid.New(),
		Name:    "Test Group",
		OwnerId: i.createdUser.Id,
	}

	groups := i.createIncomeWithGroups([]groupModel.Group{group})

	newGroup := groupModel.Group{
		Id:      uuid.New(),
		Name:    "New Test Group",
		OwnerId: i.createdUser.Id,
	}

	i.CreateEntity(newGroup)

	updatedHouse := model.UpdateIncomeRequest{
		GroupIds: []uuid.UUID{group.Id, newGroup.Id},
	}

	err := i.repository.UpdateByRequest(groups.Id, updatedHouse)

	assert.Nil(i.T(), err)

	response, err := i.repository.FindById(groups.Id)
	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []groupModel.Group{group, newGroup}, response.Groups)
}

func (i *IncomeRepositoryTestSuite) Test_UpdateByRequest_WithDeletingGroups() {
	group := groupModel.Group{
		Id:      uuid.New(),
		Name:    "Test Group",
		OwnerId: i.createdUser.Id,
	}

	house := i.createIncomeWithGroups([]groupModel.Group{group})

	updatedHouse := model.UpdateIncomeRequest{
		GroupIds: []uuid.UUID{},
	}

	err := i.repository.UpdateByRequest(house.Id, updatedHouse)

	assert.Nil(i.T(), err)

	response, err := i.repository.FindById(house.Id)
	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []groupModel.Group{}, response.Groups)
}

func (i *IncomeRepositoryTestSuite) Test_UpdateByRequest_WithMissingId() {
	assert.NotNil(i.T(), i.repository.UpdateByRequest(uuid.New(), model.UpdateIncomeRequest{}))
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

	err := i.repository.Create(&income)

	assert.Nil(i.T(), err)

	return
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
