package main

import (
	"encoding/json"
  "encoding/base64"
	"net/http"
	"strings"
  "log"
  "strconv"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
 

func getIDFromRequest(request *http.Request) string {
	params := mux.Vars(request)
	id := params["id"]
	return id
}

type dbResponse struct {
    *gorm.DB
    Cursor string
}


func decodeCursor(encodedCursor string) ([]byte){
  decoded, _ := base64.StdEncoding.DecodeString(encodedCursor)   
  return decoded
}


func GetEvents(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	defer db.Close()
	var events = []Event{}
  var results *gorm.DB

  ////
  // Experimenting with different sorting keys and directions
  ////
  // db = db.Order("updated_at desc")
  // db = db.Order("updated_at asc")
  // db = db.Order("created_at desc").Order("id")
  db = db.Order("created_at asc").Order("id asc")
  // db = db.Order("id asc")


  // Experimenting with pagination methods
  // 1) Offset + limit pagination - inefficient and prone to error
  limit := r.FormValue("limit")
  offset := r.FormValue("offset")
  if (offset != "") {
    db = db.Offset(offset)
  }
  if (limit != "") {
    db = db.Limit(limit)
  }

  // 2) Cursor (keyset, seek) pagination - efficient and reliable
  cursor := r.FormValue("cursor")
  if (cursor != "") {
    decoded := decodeCursor(cursor)

	  log.Printf("to update: %v", decoded)
    db = db.Where("created_at >= ?", decoded)
  }
	results = db.Find(&events)
  lastRecord := &events[len(events)-1]
  newCursor := lastRecord.CreatedAt.String()
  
  var dbResponse = dbResponse{
    DB: results,
    Cursor: base64.StdEncoding.EncodeToString([]byte(newCursor))}
	json.NewEncoder(w).Encode(dbResponse)
}

func GetEventCount(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	defer db.Close()
  var count uint
	db.Model(&Event{}).Count(&count)
	json.NewEncoder(w).Encode(count)
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := getIDFromRequest(r)
	db := getDB()
	defer db.Close()

	var event Event
	db.First(&event, "id = ?", id)

	s := reflect.ValueOf(&event).Elem()
	typeOfT := s.Type()

	var toUpdate = make(map[string]interface{})
	for i := 0; i < s.NumField(); i++ {
		attributeName := strings.ToLower(typeOfT.Field(i).Name)
		postValue := r.PostFormValue(attributeName)
		if postValue != "" {
			toUpdate[attributeName] = postValue
		}
	}
	log.Printf("to update: %v", toUpdate)
	db.Model(event).Updates(toUpdate)
	result := db.First(&event, "id = ?", id)
	json.NewEncoder(w).Encode(result)
}

func createFakeEvent(w http.ResponseWriter, r *http.Request) {
  count, _ := strconv.Atoi(r.FormValue("count"))
	generateFakeEvent(count)
}

func main() {
	router := mux.NewRouter()

	// Events
	router.HandleFunc("/events", GetEvents).Methods("GET")
	router.HandleFunc("/events/{id}", UpdateEvent).Methods("PUT")
	router.HandleFunc("/events_count", GetEventCount).Methods("GET")
	router.HandleFunc("/fakedata", createFakeEvent).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
