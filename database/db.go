package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBinitialize() error {
	database, err := gorm.Open(sqlite.Open("go-file-server.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	DB = database
	database.AutoMigrate(&User{})
	return nil
}
