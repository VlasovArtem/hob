package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"reflect"
)

var (
	UserRepositoryType = reflect.TypeOf(UserRepositoryObject{})
	entity             = model.User{}
)

type UserRepositoryObject struct {
	database db.ModeledDatabase
}

func NewUserRepository(service db.DatabaseService) UserRepository {
	return &UserRepositoryObject{
		db.ModeledDatabase{
			DatabaseService: service,
			Model:           entity,
		},
	}
}

func (u *UserRepositoryObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewUserRepository(factory.FindRequiredByObject(db.DatabaseObject{}).(db.DatabaseService))
}

func (u *UserRepositoryObject) GetEntity() interface{} {
	return entity
}

type UserRepository interface {
	Create(user model.User) (model.User, error)
	FindById(id uuid.UUID) (model.User, error)
	ExistsById(id uuid.UUID) bool
	ExistsByEmail(email string) bool
	FindByEmail(email string) (model.User, error)
	Update(user model.User) error
	Delete(id uuid.UUID) error
}

func (u *UserRepositoryObject) Create(user model.User) (model.User, error) {
	return user, u.database.Create(&user)
}

func (u *UserRepositoryObject) FindById(id uuid.UUID) (user model.User, err error) {
	return user, u.database.Find(&user, id)
}

func (u *UserRepositoryObject) ExistsById(id uuid.UUID) bool {
	return u.database.Exists(id)
}

func (u *UserRepositoryObject) ExistsByEmail(email string) bool {
	return u.database.ExistsBy("email = ?", email)
}

func (u *UserRepositoryObject) FindByEmail(email string) (user model.User, err error) {
	return user, u.database.FindBy(&user, "email = ?", email)
}

func (u *UserRepositoryObject) Delete(id uuid.UUID) error {
	return u.database.Delete(id)
}

func (u *UserRepositoryObject) Update(user model.User) error {
	return u.database.Update(user.Id, user)
}
