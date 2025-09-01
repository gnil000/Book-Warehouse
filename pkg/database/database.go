package database

import (
	"gin_main/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabaseConnection(config *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.Database.BookDB), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot open database connection", err)
	}
	return db
}
