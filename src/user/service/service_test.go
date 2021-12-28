package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

var (
	userRepository *mocks.UserRepository
)

func generateService() UserService {
	userRepository = new(mocks.UserRepository)

	return NewUserService(userRepository)
}

func Test_Add(t *testing.T) {
	service := generateService()

	createUserRequest := mocks.GenerateCreateUserRequest()

	var expected model.User

	userRepository.On("ExistsByEmail", createUserRequest.Email).Return(false)
	userRepository.On("Create", mock.Anything).Return(
		func(user model.User) model.User {
			expected = user
			return user
		}, nil)

	response, err := service.Add(createUserRequest)

	assert.Nil(t, err)
	assert.Equal(t, expected.ToDto(), response)
}

func Test_Add_WithExistingEmail(t *testing.T) {
	service := generateService()

	createUserRequest := mocks.GenerateCreateUserRequest()

	userRepository.On("ExistsByEmail", createUserRequest.Email).Return(true)

	response, err := service.Add(createUserRequest)

	assert.Equal(t, errors.New(fmt.Sprintf("user with '%s' already exists", createUserRequest.Email)), err)
	assert.Equal(t, model.UserDto{}, response)
	userRepository.AssertNotCalled(t, "Create")
}

func Test_Update(t *testing.T) {
	service := generateService()

	id, request := mocks.GenerateUpdateUserRequest()

	userRepository.On("ExistsById", id).Return(true)
	userRepository.On("Update", id, mock.Anything).Return(nil)

	err := service.Update(id, request)

	assert.Nil(t, err)
}

func Test_Update_WithNotExists(t *testing.T) {
	service := generateService()

	id, request := mocks.GenerateUpdateUserRequest()

	userRepository.On("ExistsById", id).Return(false)

	err := service.Update(id, request)

	assert.Equal(t, int_errors.NewErrNotFound("user with id %s not found", id), err)

	userRepository.AssertNotCalled(t, "Update", id, mock.Anything)
}

func Test_FindById(t *testing.T) {
	service := generateService()

	user := mocks.GenerateUser()

	userRepository.On("FindById", user.Id).Return(user, nil)

	response, err := service.FindById(user.Id)

	assert.Nil(t, err)
	assert.Equal(t, user.ToDto(), response)
}

func Test_FindById_WithNotExistsUser(t *testing.T) {
	service := generateService()

	id := uuid.New()

	userRepository.On("FindById", id).Return(model.User{}, gorm.ErrRecordNotFound)

	response, err := service.FindById(id)

	assert.Equal(t, int_errors.NewErrNotFound("user with id %s not found", id), err)
	assert.Equal(t, model.UserDto{}, response)
}

func Test_FindById_WithError(t *testing.T) {
	service := generateService()

	id := uuid.New()
	expectedError := errors.New("error")

	userRepository.On("FindById", id).Return(model.User{}, expectedError)

	response, err := service.FindById(id)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.UserDto{}, response)
}

func Test_ExistsById(t *testing.T) {
	service := generateService()

	id := uuid.New()

	userRepository.On("ExistsById", id).Return(true)

	assert.True(t, service.ExistsById(id))
}

func Test_ExistsById_WithoutUser(t *testing.T) {
	service := generateService()

	id := uuid.New()

	userRepository.On("ExistsById", id).Return(false)

	assert.False(t, service.ExistsById(id))
}
