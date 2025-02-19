package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB stellt sicher, dass nur eine DB-Verbindung existiert
func GetDB() *gorm.DB {
	once.Do(
		func() {
			dsn := "host=localhost user=postgres password=example  port=5555"
			var err error
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatal("Fehler beim Öffnen der DB:", err)
			}
		},
	)
	return db
}
