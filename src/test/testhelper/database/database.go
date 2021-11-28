package database

import (
	"gorm.io/gorm"
	"log"
)

func DropTable(db *gorm.DB, model interface{}) {
	err := db.Migrator().DropTable(model)

	if err != nil {
		log.Fatal(err)
	}
}

func CreateTable(db *gorm.DB, model interface{}) {
	err := db.AutoMigrate(model)

	if err != nil {
		log.Fatal(err)
	}
}

func RecreateTable(db *gorm.DB, model interface{}) {
	DropTable(db, model)
	CreateTable(db, model)
}
