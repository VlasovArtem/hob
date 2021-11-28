package respository

import (
	"db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"house/mocks"
	"house/model"
	"log"
	"test/testhelper/database"
	"testing"
	userMocks "user/mocks"
	userModel "user/model"
)

type HouseRepositoryTestSuite struct {
	suite.Suite
	createdUser userModel.User
	database    db.DatabaseService
	repository  HouseRepository
}

func (h *HouseRepositoryTestSuite) SetupSuite() {
	config := db.NewDefaultDatabaseConfiguration()
	config.DBName = "hob_test"
	h.database = db.NewDatabaseService(config)
	h.repository = NewHouseRepository(h.database)
	err := h.database.D().AutoMigrate(model.House{})

	if err != nil {
		log.Fatal(err)
	}

	h.createdUser = userMocks.GenerateUser()
	err = h.database.Create(&h.createdUser)

	if err != nil {
		log.Fatal(err)
	}
}

func TestHouseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(HouseRepositoryTestSuite))
}

func (h *HouseRepositoryTestSuite) TearDownSuite() {
	database.DropTable(h.database.D(), model.House{})
	database.DropTable(h.database.D(), userModel.User{})
}

func (h *HouseRepositoryTestSuite) Test_Create() {
	house := mocks.GenerateHouse(h.createdUser.Id)

	actual, err := h.repository.Create(house)

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), house, actual)
}

func (h *HouseRepositoryTestSuite) Test_Creat_WithMissingUser() {
	house := mocks.GenerateHouse(uuid.New())

	actual, err := h.repository.Create(house)

	assert.NotNil(h.T(), err)
	assert.Equal(h.T(), house, actual)
}

func (h *HouseRepositoryTestSuite) Test_FindById() {
	house := h.createHouse()

	actual, err := h.repository.FindResponseById(house.Id)

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), house.ToDto(), actual)
}

func (h *HouseRepositoryTestSuite) Test_FindById_WithMissingId() {
	actual, err := h.repository.FindResponseById(uuid.New())

	assert.ErrorIs(h.T(), err, gorm.ErrRecordNotFound)
	assert.Equal(h.T(), model.HouseDto{}, actual)
}

func (h *HouseRepositoryTestSuite) Test_FindByUserId() {
	house := h.createHouse()

	actual := h.repository.FindResponseByUserId(house.UserId)

	var actualResponse model.HouseDto

	for _, response := range actual {
		if response.Id == house.Id {
			actualResponse = response
			break
		}
	}
	assert.Equal(h.T(), house.ToDto(), actualResponse)
}

func (h *HouseRepositoryTestSuite) Test_FindByUserId_WithMissingId() {
	actual := h.repository.FindResponseByUserId(uuid.New())

	assert.Equal(h.T(), []model.HouseDto{}, actual)
}

func (h *HouseRepositoryTestSuite) Test_ExistsById() {
	house := h.createHouse()

	actual := h.repository.ExistsById(house.Id)

	assert.True(h.T(), actual)
}

func (h *HouseRepositoryTestSuite) Test_ExistsById_WithMissingId() {
	actual := h.repository.ExistsById(uuid.New())

	assert.False(h.T(), actual)
}

func (h *HouseRepositoryTestSuite) createHouse() model.House {
	house := mocks.GenerateHouse(h.createdUser.Id)

	create, err := h.repository.Create(house)

	assert.Nil(h.T(), err)

	return create
}
