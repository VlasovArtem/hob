package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/VlasovArtem/hob/src/user/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserServiceStr struct {
	repository repository.UserRepository
}

func (u *UserServiceStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(repository.UserRepositoryObject{}),
	}
}

func (u *UserServiceStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewUserService(dependency.FindRequiredDependency[repository.UserRepositoryObject, repository.UserRepository](factory))
}

func NewUserService(repository repository.UserRepository) UserService {
	return &UserServiceStr{repository}
}

type UserService interface {
	transactional.Transactional[UserService]
	Add(request model.CreateUserRequest) (model.UserDto, error)
	Update(id uuid.UUID, request model.UpdateUserRequest) error
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) (model.UserDto, error)
	ExistsById(id uuid.UUID) bool
	VerifyUser(email string, password string) (model.UserDto, error)
	Transactional(db *gorm.DB) UserService
}

func (u *UserServiceStr) Add(request model.CreateUserRequest) (response model.UserDto, err error) {
	if u.repository.ExistsByEmail(request.Email) {
		return response, errors.New(fmt.Sprintf("user with '%s' already exists", request.Email))
	}
	if request.Email == "" {
		return response, errors.New(fmt.Sprintf("email is missing"))
	}
	if request.Password == "" {
		return response, errors.New(fmt.Sprintf("password is missing"))
	}

	entity := request.ToEntity()
	if err := u.repository.Create(&entity); err != nil {
		return response, err
	} else {
		return entity.ToDto(), err
	}
}

func (u *UserServiceStr) Update(id uuid.UUID, request model.UpdateUserRequest) error {
	if !u.ExistsById(id) {
		return int_errors.NewErrNotFound("user with id %s not found", id)
	}
	return u.repository.Update(id, request)
}

func (u *UserServiceStr) Delete(id uuid.UUID) error {
	return u.repository.Delete(id)
}

func (u *UserServiceStr) FindById(id uuid.UUID) (response model.UserDto, err error) {
	if user, err := u.repository.First(id); err != nil {
		return response, database.HandlerFindError(err, "user with id %s not found", id)
	} else {
		return user.ToDto(), err
	}
}

func (u *UserServiceStr) ExistsById(id uuid.UUID) bool {
	return u.repository.Exists(id)
}

func (u *UserServiceStr) VerifyUser(email string, password string) (response model.UserDto, err error) {
	if !u.repository.Verify(email, []byte(password)) {
		return response, errors.New("credentials are not valid")
	}

	if user, err := u.repository.FindByEmail(email); err != nil {
		return response, err
	} else {
		return user.ToDto(), err
	}
}

func (u *UserServiceStr) Transactional(db *gorm.DB) UserService {
	return &UserServiceStr{
		&repository.UserRepositoryObject{
			ModeledDatabase: u.repository.Transactional(db),
		},
	}
}
