package db

import (
	"log"
	"mint-redeem-workflow/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	var err error
	Db, err = gorm.Open(sqlite.Open("mint-redeem.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	Db.AutoMigrate(&models.Request{})
}
