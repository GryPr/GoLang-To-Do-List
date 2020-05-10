package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// ToDoItem struct contains the data of a to-do item
type ToDoItem struct {
	ID          int    `gorm:"primary_key;auto_increment"`
	Description string `json:"description"`
	Completion  bool   `json:"completion"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

// Health returns the API health
func Health(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")                       // Logs to console API Health
	w.Header().Set("Content-Type", "application/json") // Sets the content type to JSON
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, `{"alive": true}`) // JSON response to client
}

// CreateItem adds a new To-Do item to the database
func CreateItem(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var todo ToDoItem
	err := decoder.Decode(&todo)
	if err != nil {
		panic(err)
	}
	log.WithFields(log.Fields{"description": todo.Description}).Info("Adding new item and saving to the database") // Logs to console the description
	db.Create(&todo)                                                                                               // Inserts the struct into the database
	result := db.Last(&todo)                                                                                       // Gets last record ordered by primary key
	w.Header().Set("Content-Type", "application/json")                                                             // Adds json header
	json.NewEncoder(w).Encode(&result.Value)                                                                       // Responds with the last record
}

// GetCompleteItems returns all the complete items
func GetCompleteItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting completed items")
	completedItems := getItemsByCompletion(true)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(&completedItems)
}

// GetIncompleteItems returns all the incomplete items
func GetIncompleteItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting incomplete items")
	incompleItems := getItemsByCompletion(false)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(&incompleItems)
}

// GetAllItems returns all the items
func GetAllItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting all items")
	var tditems []ToDoItem // Array of ToDoItem struct
	allItems := db.Order("id").Find(&tditems).Value
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(&allItems)
}

// UpdateCompletionItem updates the completion of an item
func UpdateCompletionItem(w http.ResponseWriter, r *http.Request) {
	log.Info("Updating specific item")
	vars := mux.Vars(r) // Gets the variable from the request
	id, _ := strconv.Atoi(vars["id"])
	decoder := json.NewDecoder(r.Body)
	var todo ToDoItem
	err := decoder.Decode(&todo)
	var td ToDoItem
	td.Completion = todo.Completion
	if err != nil {
		panic(err)
	}
	foundItems := getItemsByID(id)
	if foundItems == false {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.WriteString(w, `{"deleted": false, "error": "Record Not Found"}`)
	} else {
		log.WithFields(log.Fields{"Id": id, "Completed": todo.Completion}).Info("Updating Item")
		db.Save(&td)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(&todo)
	}

}

// DeleteItem receives an ID from the request and deletes an item
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	// Gets the ID from the request and converts it from string to int
	vars := mux.Vars(r) // Gets the variable from the request
	id, _ := strconv.Atoi(vars["id"])

	err := getItemsByID(id)
	if err == false {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.WriteString(w, `{"deleted": false, "error" "Record Not Found"}`) // Error JSON
	} else {
		log.WithFields(log.Fields{"Id": id}).Info("Deleting Item")
		todo := &ToDoItem{}
		db.First(&todo, id)
		db.Delete(&todo)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.WriteString(w, `{"deleted": true}`)
	}

}

func getItemsByCompletion(completion bool) interface{} {
	var tditems []ToDoItem                                                   // Array of ToDoItem struct
	toDoItems := db.Where("completion = ?", completion).Find(&tditems).Value // Finding which database items are
	return toDoItems
}

// getItemsByID checks for the existence of an item by a specific ID
func getItemsByID(id int) bool {
	todo := &ToDoItem{}
	result := db.First(&todo, id)
	if result.Error != nil {
		log.Warn("Item not found in the database")
		return false
	}
	return true
}
