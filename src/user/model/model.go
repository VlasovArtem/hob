package model

import (
	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `gorm:"primarykey"`
	FirstName string
	LastName  string
	Password  []byte
	Email     string `gorm:"unique"`
}

type UserDto struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Email     string
}

type CreateUserRequest struct {
	FirstName string
	LastName  string
	Password  string
	Email     string
}

type UpdateUserRequest struct {
	FirstName string
	LastName  string
	Password  string
}

func (u CreateUserRequest) ToEntity() User {
	return User{
		Id:        uuid.New(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Password:  []byte(u.Password),
		Email:     u.Email,
	}
}

func (u User) ToDto() UserDto {
	return UserDto{
		Id:        u.Id,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}
