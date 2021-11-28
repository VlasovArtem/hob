package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
)

type UserRepositoryObject struct {
	database db.DatabaseService
}

func NewUserRepository(service db.DatabaseService) UserRepository {
	return &UserRepositoryObject{service}
}

func (u *UserRepositoryObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewUserRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService)),
	)
}

func (u *UserRepositoryObject) GetEntity() interface{} {
	return model.User{}
}

type UserRepository interface {
	Create(user model.User) (model.User, error)
	FindById(id uuid.UUID) (model.User, error)
	ExistsById(id uuid.UUID) bool
	ExistsByEmail(email string) bool
}

func (u *UserRepositoryObject) Create(user model.User) (model.User, error) {
	return user, u.database.Create(&user)
}

func (u *UserRepositoryObject) FindById(id uuid.UUID) (model.User, error) {
	user := model.User{}

	return user, u.database.FindById(&user, id)
}

func (u *UserRepositoryObject) ExistsById(id uuid.UUID) bool {
	return u.database.ExistsById(model.User{}, id)
}

func (u *UserRepositoryObject) ExistsByEmail(email string) bool {
	return u.database.ExistsByQuery(model.User{}, "email = ?", email)
}

func (u *UserRepositoryObject) Migrate() interface{} {
	return model.User{}
}
