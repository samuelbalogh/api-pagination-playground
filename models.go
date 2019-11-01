package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/icrowley/fake"
)


type Event struct {
	ID           uuid.UUID `gorm:"primary_key; type:uuid; colum:id;default:uuid_generate_v4()`
	Title        string
	Starts       time.Time
	Ends         time.Time
	Description  string
	Category     string
	IsRecurring  bool
	Frequency    string
	CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time
}


func (event *Event) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.New()
	scope.SetColumn("ID", id)
	return nil
}


func generateFakeEvent(count int) (){
	db := getDB()
	defer db.Close()
	for i := 0; i<count; i++ {
		title := fake.Sentence()
		starts := time.Now()
		ends := time.Now()
		description := fmt.Sprintf("Description %d", i)
		event := Event{
			Title: title, 
			Starts: starts, 
			Ends: ends,
			Description: description,
		}
		db.Create(&event)
		fmt.Printf("%+v \n", event)
	}
}
