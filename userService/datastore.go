package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

//CreateDBConnection -creates a conection to the sql database
func CreateDBConnection() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	return gorm.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, password),
	)

}
