package mocks

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
)

func GenerateCreateUserRequest() model.CreateUserRequest {
	return model.CreateUserRequest{
		FirstName: "First Name",
		LastName:  "Last Name",
		Password:  "password",
		Email:     "mail@mail.com",
	}
}

func GenerateUpdateUserRequest() (uuid.UUID, model.UpdateUserRequest) {
	return uuid.New(), model.UpdateUserRequest{
		FirstName: "First Name New",
		LastName:  "Last Name New",
		Password:  "password-new",
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

func GenerateUserResponse() model.UserDto {
	return model.UserDto{
		Id:        uuid.New(),
		FirstName: "First Name",
		LastName:  "Last Name",
		Email:     "mail@mail.com",
	}
}
