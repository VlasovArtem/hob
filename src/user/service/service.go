package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/VlasovArtem/hob/src/user/repository"
	"github.com/google/uuid"
	"reflect"
)

var UserServiceType = reflect.TypeOf(UserServiceObject{})

type UserServiceObject struct {
	repository repository.UserRepository
}

func (u *UserServiceObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewUserService(factory.FindRequiredByObject(repository.UserRepositoryObject{}).(repository.UserRepository))
}

func NewUserService(repository repository.UserRepository) UserService {
	return &UserServiceObject{repository}
}

type UserService interface {
	Add(request model.CreateUserRequest) (model.UserDto, error)
	Update(id uuid.UUID, request model.UpdateUserRequest) error
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) (model.UserDto, error)
	ExistsById(id uuid.UUID) bool
	VerifyUser(email string, password string) (model.UserDto, error)
}

func (u *UserServiceObject) Add(request model.CreateUserRequest) (response model.UserDto, err error) {
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
		return user.ToDto(), err
	}
}

func (u *UserServiceObject) Update(id uuid.UUID, request model.UpdateUserRequest) error {
	if !u.ExistsById(id) {
		return int_errors.NewErrNotFound("user with id %s not found", id)
	}
	return u.repository.Update(id, request)
}

func (u *UserServiceObject) Delete(id uuid.UUID) error {
	return u.repository.Delete(id)
}

func (u *UserServiceObject) FindById(id uuid.UUID) (response model.UserDto, err error) {
	if user, err := u.repository.FindById(id); err != nil {
		return response, database.HandlerFindError(err, "user with id %s not found", id)
	} else {
		return user.ToDto(), err
	}
}

func (u *UserServiceObject) ExistsById(id uuid.UUID) bool {
	return u.repository.ExistsById(id)
}

func (u *UserServiceObject) VerifyUser(email string, password string) (response model.UserDto, err error) {
	if user, err := u.repository.FindByEmail(email); err != nil {
		return response, database.HandlerFindError(err, "user not found")
	} else if string(user.Password) != password {
		return response, errors.New("credentials is not valid")
	} else {
		return user.ToDto(), nil
	}
}
