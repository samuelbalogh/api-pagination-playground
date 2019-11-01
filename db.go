package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)


func getDB() (*gorm.DB) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=calendar dbname=postgres password=calendar sslmode=disable")

	checkErr(err)
	return db
}
