package database

import (
	"errors"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	"gorm.io/gorm"
)

func HandlerFindError(err error, message string, args ...interface{}) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return int_errors.NewErrNotFound(message, args...)
	}
	return err
}
