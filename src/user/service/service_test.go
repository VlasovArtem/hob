package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type UserServiceTestSuite struct {
	testhelper.MockTestSuite[UserService]
	userRepository *mocks.UserRepository
}

func TestUserServiceTestSuite(t *testing.T) {
	ts := &UserServiceTestSuite{}
	ts.TestObjectGenerator = func() UserService {
		ts.userRepository = new(mocks.UserRepository)

		return NewUserService(ts.userRepository)
	}

	suite.Run(t, ts)
}

func (u *UserServiceTestSuite) Test_Add() {
	createUserRequest := mocks.GenerateCreateUserRequest()

	var expected model.User

	u.userRepository.On("ExistsByEmail", createUserRequest.Email).Return(false)
	u.userRepository.On("Create", mock.Anything).Return(
		func(user model.User) model.User {
			expected = user
			return user
		}, nil)

	response, err := u.TestO.Add(createUserRequest)

	assert.Nil(u.T(), err)
	assert.Equal(u.T(), expected.ToDto(), response)
}

func (u *UserServiceTestSuite) Test_Add_WithExistingEmail() {
	createUserRequest := mocks.GenerateCreateUserRequest()

	u.userRepository.On("ExistsByEmail", createUserRequest.Email).Return(true)

	response, err := u.TestO.Add(createUserRequest)

	assert.Equal(u.T(), errors.New(fmt.Sprintf("user with '%s' already exists", createUserRequest.Email)), err)
	assert.Equal(u.T(), model.UserDto{}, response)
	u.userRepository.AssertNotCalled(u.T(), "Create")
}

func (u *UserServiceTestSuite) Test_Update() {
	id, request := mocks.GenerateUpdateUserRequest()

	u.userRepository.On("ExistsById", id).Return(true)
	u.userRepository.On("Update", id, mock.Anything).Return(nil)

	err := u.TestO.Update(id, request)

	assert.Nil(u.T(), err)
}

func (u *UserServiceTestSuite) Test_Update_WithNotExists() {
	id, request := mocks.GenerateUpdateUserRequest()

	u.userRepository.On("ExistsById", id).Return(false)

	err := u.TestO.Update(id, request)

	assert.Equal(u.T(), int_errors.NewErrNotFound("user with id %s not found", id), err)

	u.userRepository.AssertNotCalled(u.T(), "Update", id, mock.Anything)
}

func (u *UserServiceTestSuite) Test_FindById() {
	user := mocks.GenerateUser()

	u.userRepository.On("FindById", user.Id).Return(user, nil)

	response, err := u.TestO.FindById(user.Id)

	assert.Nil(u.T(), err)
	assert.Equal(u.T(), user.ToDto(), response)
}

func (u *UserServiceTestSuite) Test_FindById_WithNotExistsUser() {
	id := uuid.New()

	u.userRepository.On("FindById", id).Return(model.User{}, gorm.ErrRecordNotFound)

	response, err := u.TestO.FindById(id)

	assert.Equal(u.T(), int_errors.NewErrNotFound("user with id %s not found", id), err)
	assert.Equal(u.T(), model.UserDto{}, response)
}

func (u *UserServiceTestSuite) Test_FindById_WithError() {
	id := uuid.New()
	expectedError := errors.New("error")

	u.userRepository.On("FindById", id).Return(model.User{}, expectedError)

	response, err := u.TestO.FindById(id)

	assert.Equal(u.T(), expectedError, err)
	assert.Equal(u.T(), model.UserDto{}, response)
}

func (u *UserServiceTestSuite) Test_ExistsById() {
	id := uuid.New()

	u.userRepository.On("ExistsById", id).Return(true)

	assert.True(u.T(), u.TestO.ExistsById(id))
}

func (u *UserServiceTestSuite) Test_ExistsById_WithoutUser() {
	id := uuid.New()

	u.userRepository.On("ExistsById", id).Return(false)

	assert.False(u.T(), u.TestO.ExistsById(id))
}
