package mocks

import (
	"fmt"
	"github.com/google/uuid"
	"user/model"
)

func GenerateCreateUserRequest() model.CreateUserRequest {
	return model.CreateUserRequest{
		FirstName: "First Name",
		LastName:  "Last Name",
		Password:  "password",
		Email:     "mail@mai.com",
	}
}

func GenerateUser() model.User {
	id := uuid.New()
	return model.User{
		Id:        id,
		FirstName: "First Name",
		LastName:  "Last Name",
		Password:  []byte("password"),
		Email:     fmt.Sprintf("mail%s@mail.com", id),
	}
}

func GenerateUserResponse() model.UserResponse {
	return model.UserResponse{
		Id:        uuid.New(),
		FirstName: "First Name",
		LastName:  "Last Name",
		Email:     "mail@mai.com",
	}
}
