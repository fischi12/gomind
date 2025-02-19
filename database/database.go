package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func GetDB() *gorm.DB {

	dsn := "host=192.168.2.188 user=postgres password=example  port=5555"
	var db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Fehler beim Öffnen der DB:", err)
	}

	return db
}
