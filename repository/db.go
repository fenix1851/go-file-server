package repository

import (
	"fmt"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	onceDB sync.Once
)

func DBinitialize() (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open("go-file-server.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	database.AutoMigrate(&User{})

	return database, nil
}

func GetDBInstance() (*gorm.DB, error) {
	var err error
	onceDB.Do(func() {
		DB, err = DBinitialize()
	})
	return DB, err
}
