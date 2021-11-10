package model

import "github.com/google/uuid"

type User struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	password  []byte
	Email     string
}

type UserResponse struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Email     string
}

type CreateUserRequest struct {
	FirstName string
	LastName  string
	Password  []byte
	Email     string
}

func (u CreateUserRequest) ToEntity() User {
	return User{
		Id:        uuid.New(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		password:  u.Password,
		Email:     u.Email,
	}
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		Id:        u.Id,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}
