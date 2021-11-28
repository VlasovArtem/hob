package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/VlasovArtem/hob/src/user/repository"
	"github.com/google/uuid"
)

type UserServiceObject struct {
	repository repository.UserRepository
}

func (u *UserServiceObject) Initialize(factory dependency.DependenciesFactory) {
	userRepository := factory.FindRequiredByObject(repository.UserRepositoryObject{}).(repository.UserRepository)

	factory.Add(NewUserService(userRepository))
}

func NewUserService(repository repository.UserRepository) UserService {
	return &UserServiceObject{repository}
}

type UserService interface {
	Add(request model.CreateUserRequest) (model.UserResponse, error)
	FindById(id uuid.UUID) (model.UserResponse, error)
	ExistsById(id uuid.UUID) bool
}

func (u *UserServiceObject) Add(request model.CreateUserRequest) (response model.UserResponse, err error) {
	if u.repository.ExistsByEmail(request.Email) {
		return response, errors.New(fmt.Sprintf("user with '%s' already exists", request.Email))
	}
	if request.Email == "" {
		return response, errors.New(fmt.Sprintf("email is missing"))
	}
	if request.Password == "" {
		return response, errors.New(fmt.Sprintf("password is missing"))
	}

	if user, err := u.repository.Create(request.ToEntity()); err != nil {
		return response, err
	} else {
		return user.ToResponse(), err
	}
}

func (u *UserServiceObject) FindById(id uuid.UUID) (response model.UserResponse, err error) {
	if user, err := u.repository.FindById(id); err != nil {
		return response, database.HandlerFindError(err, fmt.Sprintf("user with id %s in not exists", id))
	} else {
		return user.ToResponse(), nil
	}
}

func (u *UserServiceObject) ExistsById(id uuid.UUID) bool {
	return u.repository.ExistsById(id)
}
