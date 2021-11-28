package database

import (
	"errors"
	"gorm.io/gorm"
)

func HandlerFindError(err error, notFoundMessage string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(notFoundMessage)
	}
	return err
}
