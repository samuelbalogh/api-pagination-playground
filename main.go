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
    Count uint
}


func decodeCursor(encodedCursor string) (string){
  decoded, _ := base64.StdEncoding.DecodeString(encodedCursor)   
  return string(decoded)
}


func GetEvents(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	defer db.Close()
	var events = []Event{}
  var results *gorm.DB

  // 1) Offset + limit pagination - inefficient and prone to error
  offset := r.FormValue("offset")
  if (offset != "") {
    db = db.Offset(offset)
  }

  // 2) Cursor (keyset, seek) pagination - efficient and reliable, 
  // but pushes some logic to the client, and does not support flexible ordering
  cursor := r.FormValue("cursor")
  if (cursor != "") {
    // The cursor in this case is a base64-encoded composite key, made up of the
    // 1) created_at and 2) id of the last record. 
    // This is to ensure that  we have unique cursors (there can be records with the same created_at)
    // NOTE: this implementation exposes the internal UUID of records in an encoded form - 
    // this is not a good practice and ideally it would be encrypted before sending to the client,
    // and decrypted when received.
    decoded := decodeCursor(cursor)
    createdAt := strings.Split(decoded, "#")[0]
    id := strings.Split(decoded, "#")[1]
    db = db.Where("created_at >= ? and (id > ? or created_at > ?)", createdAt, id, createdAt)

    // We have to ensure consistent ordering for this to work
    db = db.Order("created_at asc").Order("id asc")
  }

  if (offset == "" && cursor == "") { 
    db = db.Order("created_at asc").Order("id asc") 
  }

  db.LogMode(true)
  orderby := r.FormValue("orderby")
  if (orderby != "") {
    // WARNING: don't do this in production ever - don't accept arbitrary values and pass them to SQL as here.
    // This is used only for making experimentation with ordering easier.
    db = db.Order(orderby)
  } 

  limit := r.FormValue("limit")
  if (limit != "") {
    db = db.Limit(limit)
  } else {
    db = db.Limit(10000)
  }

  count := GetEventCount()

	results = db.Find(&events)
  var dbResp dbResponse
  if (len(events) == 0) {
    dbResp = dbResponse{DB: results, Cursor: "", Count: count}
  } else {
    lastRecord := &events[len(events)-1]
    lastCreatedAt := lastRecord.CreatedAt.Format("2006-01-02T15:04:05.000000")
    lastId := lastRecord.ID
    newCursor := lastCreatedAt + "#" + lastId.String()
  
    dbResp = dbResponse{
      DB: results,
      Cursor: base64.StdEncoding.EncodeToString([]byte(newCursor)),
      Count: count}
  }
	json.NewEncoder(w).Encode(dbResp)
}

func GetEventCount() uint {
	db := getDB()
	defer db.Close()
  var count uint
	db.Model(&Event{}).Count(&count)
  return count
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
	router.HandleFunc("/fakedata", createFakeEvent).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
