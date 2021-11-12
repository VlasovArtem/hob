package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"user/model"
)

type userServiceObject struct {
	users map[uuid.UUID]model.User
}

func NewUserService() UserService {
	return &userServiceObject{
		users: make(map[uuid.UUID]model.User),
	}
}

type UserService interface {
	Add(request model.CreateUserRequest) (model.UserResponse, error)
	FindById(id uuid.UUID) (model.UserResponse, error)
	ExistsById(id uuid.UUID) bool
	ExistsByEmail(email string) bool
}

func (u *userServiceObject) Add(request model.CreateUserRequest) (model.UserResponse, error) {
	if u.ExistsByEmail(request.Email) {
		return model.UserResponse{}, errors.New(fmt.Sprintf("user with '%s' already exists", request.Email))
	}

	user := request.ToEntity()

	u.users[user.Id] = user

	return user.ToResponse(), nil
}

func (u *userServiceObject) FindById(id uuid.UUID) (model.UserResponse, error) {
	if user, ok := u.users[id]; ok {
		return user.ToResponse(), nil
	}
	return model.UserResponse{}, errors.New(fmt.Sprintf("user with %s is not found", id))
}

func (u *userServiceObject) ExistsById(id uuid.UUID) bool {
	_, ok := u.users[id]
	return ok
}

func (u *userServiceObject) ExistsByEmail(email string) bool {
	for _, user := range u.users {
		if user.Email == email {
			return true
		}
	}
	return false
}
