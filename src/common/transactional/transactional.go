package transactional

import "gorm.io/gorm"

type Transactional[T any] interface {
	Transactional(tx *gorm.DB) T
}
