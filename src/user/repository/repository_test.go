package repository

import (
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	"github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type UserRepositoryTestSuite struct {
	database.DBTestSuite
	repository UserRepository
}

func (p *UserRepositoryTestSuite) SetupSuite() {
	p.InitDBTestSuite()

	p.CreateRepository(
		func(service db.DatabaseService) {
			p.repository = NewUserRepository(service)
		},
	).
		AddMigrators(model.User{})
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
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
	user := p.createUser()

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
	user := p.createUser()

	assert.True(p.T(), p.repository.ExistsById(user.Id))
}

func (p *UserRepositoryTestSuite) Test_ExistsById_WithNotExistsUser() {
	id := uuid.New()

	assert.False(p.T(), p.repository.ExistsById(id))
}

func (p *UserRepositoryTestSuite) Test_ExistsByEmail() {
	user := p.createUser()

	assert.True(p.T(), p.repository.ExistsByEmail(user.Email))
}

func (p *UserRepositoryTestSuite) Test_ExistsByEmail_WithNotExistsUser() {
	assert.False(p.T(), p.repository.ExistsByEmail("email@mail.com"))
}

func (p *UserRepositoryTestSuite) Test_Delete() {
	provider := p.createUser()

	err := p.repository.Delete(provider.Id)

	assert.Nil(p.T(), err)
	assert.False(p.T(), p.repository.ExistsById(provider.Id))
}

func (p *UserRepositoryTestSuite) Test_Update() {
	user := p.createUser()

	updated := model.UpdateUserRequest{
		FirstName: "New First Name",
		LastName:  "New Last Name",
		Password:  "new",
	}
	err := p.repository.Update(user.Id, updated)

	assert.Nil(p.T(), err)
	user1, err := p.repository.FindById(user.Id)
	assert.Nil(p.T(), err)
	assert.Equal(p.T(), model.User{
		Id:        user.Id,
		FirstName: "New First Name",
		LastName:  "New Last Name",
		Password:  []byte("new"),
		Email:     user.Email,
	}, user1)
}

func (p *UserRepositoryTestSuite) createUser() model.User {
	user := mocks.GenerateUser()

	p.CreateEntity(&user)

	return user
}
