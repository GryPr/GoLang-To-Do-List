package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "gryphticon"
	dbname   = "todo"
)

var psqlInfo string = fmt.Sprintf("host=%s port =%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
var db, err = gorm.Open("postgres", psqlInfo)

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
func health(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")                       // Logs to console API Health
	w.Header().Set("Content-Type", "application/json") // Sets the content type to JSON
	io.WriteString(w, `{"alive": true}`)               // JSON response to client
}

// createItem adds a new To-Do item to the database
func createItem(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(result.Value)                                                                        // Responds with the last record
}

// Gets all the complete items
func getCompleteItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting completed items")
	completedItems := getItemsByCompletion(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completedItems)
}

// Gets all the incomplete items
func getIncompleteItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting incomplete items")
	incompleItems := getItemsByCompletion(false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incompleItems)
}

// Updates the completion of an item
func updateItem(w http.ResponseWriter, r *http.Request) {
	log.Info("Updating specific item")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"]) // string to int

	err := getItemsByID(id)
	if err == false {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": false, "error": "Record Not Found"}`)
	} else {
		completed, _ := strconv.ParseBool(r.FormValue("completed")) // Parses the bool from the POST
		log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating Item")
		todo := &ToDoItem{}
		db.First(&todo, id)
		todo.Completion = completed
		w.Header().Set("Content-Type", "application/json")
	}

}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	// Gets the ID from the request and converts it from string to int
	vars := mux.Vars(r) // Gets the variable from the request
	id, _ := strconv.Atoi(vars["id"])

	err := getItemsByID(id)
	if err == false {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": false, "error" "Record Not Found"}`) // Error JSON
	} else {
		log.WithFields(log.Fields{"Id": id}).Info("Deleting Item")
		todo := &ToDoItem{}
		db.First(&todo, id)
		db.Delete(&todo)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": true}`)
	}

}

func getItemsByCompletion(completion bool) interface{} {
	var tditems []ToDoItem
	toDoItems := db.Where("completion = ?", completion).Find(&tditems).Value
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

func main() {
	defer db.Close()

	// Setting up database

	if err != nil {
		panic(err)
	}
	if err == nil {
		fmt.Println("Sucessfully connected to the database")
	}

	db.Debug().DropTableIfExists(&ToDoItem{})
	db.Debug().AutoMigrate(&ToDoItem{})

	log.Info("Starting API Server")

	// Setting up the router
	router := mux.NewRouter()
	router.HandleFunc("/ping", health).Methods("GET")
	router.HandleFunc("/todo", createItem).Methods("POST")
	router.HandleFunc("/todo-complete", getCompleteItems).Methods("GET")
	router.HandleFunc("/todo-incomplete", getIncompleteItems).Methods("GET")
	router.HandleFunc("/todo/{id}", updateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", deleteItem).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}
