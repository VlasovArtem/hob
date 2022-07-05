package repository

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/transactional"
	"github.com/VlasovArtem/hob/src/db"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryObject struct {
	db.ModeledDatabase[model.User]
}

func NewUserRepository(service db.DatabaseService) UserRepository {
	return &UserRepositoryObject{
		ModeledDatabase: db.NewModeledDatabase(model.User{}, service),
	}
}

func (u *UserRepositoryObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewUserRepository(factory.FindRequiredByObject(db.Database{}).(db.DatabaseService))
}

type UserRepository interface {
	db.ModeledDatabase[model.User]
	transactional.Transactional[UserRepository]
	FindByEmail(email string) (model.User, error)
	ExistsByEmail(email string) bool
	Verify(email string, password []byte) bool
	UpdateByRequest(id uuid.UUID, user model.UpdateUserRequest) error
}

func (u *UserRepositoryObject) FindByEmail(email string) (user model.User, error error) {
	return u.FirstBy("email = ?", email)
}

func (u *UserRepositoryObject) ExistsByEmail(email string) bool {
	return u.ExistsBy("email = ?", email)
}

func (u *UserRepositoryObject) Verify(email string, password []byte) bool {
	return u.ExistsBy("email = ? AND password = ?", email, password)
}

func (u *UserRepositoryObject) UpdateByRequest(id uuid.UUID, user model.UpdateUserRequest) error {
	return u.Update(id, struct {
		FirstName string
		LastName  string
		Password  []byte
	}{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  []byte(user.Password),
	})
}

func (u *UserRepositoryObject) Transactional(tx *gorm.DB) UserRepository {
	return &UserRepositoryObject{
		ModeledDatabase: db.NewTransactionalModeledDatabase(model.User{}, tx),
	}
}
