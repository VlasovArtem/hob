package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"user/model"
)

func Test_user_AddUser(t *testing.T) {
	type fields struct {
		users map[uuid.UUID]model.User
	}
	type args struct {
		request model.CreateUserRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		response func(id uuid.UUID) model.UserResponse
		err      error
	}{
		{
			name: "save",
			fields: fields{
				users: make(map[uuid.UUID]model.User),
			},
			args: args{
				request: model.CreateUserRequest{
					FirstName: "First Name",
					LastName:  "Last Name",
					Password:  []byte("password"),
					Email:     "mail@mail.com",
				},
			},
			response: func(id uuid.UUID) model.UserResponse {
				return model.UserResponse{
					Id:        id,
					FirstName: "First Name",
					LastName:  "Last Name",
					Email:     "mail@mail.com",
				}
			},
			err: nil,
		},
		{
			name: "with existing email",
			fields: fields{
				users: func() map[uuid.UUID]model.User {
					m := make(map[uuid.UUID]model.User)

					user := model.CreateUserRequest{
						FirstName: "First Name",
						LastName:  "Last Name",
						Password:  []byte("password"),
						Email:     "mail@mail.com",
					}.ToEntity()

					m[user.Id] = user

					return m
				}(),
			},
			args: args{
				request: model.CreateUserRequest{
					FirstName: "First Name",
					LastName:  "Last Name",
					Password:  []byte("password"),
					Email:     "mail@mail.com",
				},
			},
			response: func(id uuid.UUID) model.UserResponse {
				return model.UserResponse{}
			},
			err: errors.New("user with 'mail@mail.com' already exists"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userServiceObject{
				users: tt.fields.users,
			}

			newUser, err := u.AddUser(tt.args.request)

			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.response(newUser.Id), newUser)
		})
	}
}

func Test_user_FindById(t *testing.T) {
	type fields struct {
		users map[uuid.UUID]model.User
	}
	type args struct {
		id uuid.UUID
	}

	f := fields{
		users: make(map[uuid.UUID]model.User),
	}

	u := model.User{
		Id:        uuid.New(),
		FirstName: "First name",
		LastName:  "Last name",
		Email:     "mail@mail.com",
	}

	f.users[u.Id] = u

	notExistingId := uuid.New()

	tests := []struct {
		name     string
		fields   fields
		args     args
		err      error
		response model.UserResponse
	}{
		{
			name:     "with existing user",
			fields:   f,
			args:     args{u.Id},
			err:      nil,
			response: u.ToResponse(),
		},
		{
			name:     "with not existing user",
			fields:   f,
			args:     args{notExistingId},
			err:      errors.New(fmt.Sprintf("user with %s is not found", notExistingId)),
			response: model.UserResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userServiceObject{
				users: tt.fields.users,
			}
			got, got1 := u.FindById(tt.args.id)
			assert.Equalf(t, tt.err, got1, "FindById(%v)", tt.args.id)
			assert.Equalf(t, tt.response, got, "FindById(%v)", tt.args.id)
		})
	}
}
