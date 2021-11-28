package repository

import (
	"db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"log"
	"test/testhelper/database"
	"testing"
	"user/mocks"
	"user/model"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	database   db.DatabaseService
	repository UserRepository
}

func (p *UserRepositoryTestSuite) SetupSuite() {
	config := db.NewDefaultDatabaseConfiguration()
	config.DBName = "hob_test"
	p.database = db.NewDatabaseService(config)
	p.repository = NewUserRepository(p.database)
	err := p.database.D().AutoMigrate(model.User{})

	if err != nil {
		log.Fatal(err)
	}
}

func (p *UserRepositoryTestSuite) TearDownSuite() {
	database.DropTable(p.database.D(), model.User{})
}

func (p *UserRepositoryTestSuite) Test_Create() {
	user := mocks.GenerateUser()

	create, err := p.repository.Create(user)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), user, create)
}

func (p *UserRepositoryTestSuite) Test_CreateWithMatchingEmail() {
	expected := mocks.GenerateUser()

	actual, err := p.repository.Create(expected)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), expected, actual)

	actual, err = p.repository.Create(expected)

	assert.Equal(p.T(), expected, actual)
	assert.NotNil(p.T(), err)
}

func (p *UserRepositoryTestSuite) Test_FindById() {
	user := p.CreateUser()

	actual, err := p.repository.FindById(user.Id)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), user, actual)
}

func (p *UserRepositoryTestSuite) Test_FindById_WithNotExistsUser() {
	id := uuid.New()

	actual, err := p.repository.FindById(id)

	assert.Equal(p.T(), gorm.ErrRecordNotFound, err)
	assert.Equal(p.T(), model.User{}, actual)
}

func (p *UserRepositoryTestSuite) Test_ExistsById() {
	user := p.CreateUser()

	assert.True(p.T(), p.repository.ExistsById(user.Id))
}

func (p *UserRepositoryTestSuite) Test_ExistsById_WithNotExistsUser() {
	id := uuid.New()

	assert.False(p.T(), p.repository.ExistsById(id))
}

func (p *UserRepositoryTestSuite) Test_ExistsByEmail() {
	user := p.CreateUser()

	assert.True(p.T(), p.repository.ExistsByEmail(user.Email))
}

func (p *UserRepositoryTestSuite) Test_ExistsByEmail_WithNotExistsUser() {
	assert.False(p.T(), p.repository.ExistsByEmail("email@mail.com"))
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (p *UserRepositoryTestSuite) CreateUser() model.User {
	user := mocks.GenerateUser()

	savedUser, err := p.repository.Create(user)

	assert.Nil(p.T(), err)

	return savedUser
}
