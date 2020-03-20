package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

var db *gorm.DB
var err error

// ToDoItem struct contains the data of a to-do item
type ToDoItem struct {
	id          int `gorm:"primary_key"`
	description string
	completion  bool
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
	description := r.FormValue("description")                                                                 // Obtaining the value from POST
	log.WithFields(log.Fields{"description": description}).Info("Adding new item and saving to the database") // Logs to console the description
	todo := &ToDoItem{description: description, completion: false}                                            // Passes ToDoItem struct by reference
	db.Create(&todo)                                                                                          // Inserts the struct into the database
	result := db.Last(&todo)                                                                                  // Gets last record ordered by primary key
	w.Header().Set("Content-Type", "application/json")                                                        // Adds json header
	json.NewEncoder(w).Encode(result.Value)                                                                   // Responds with the last record
}

func main() {
	defer db.Close()

	// Setting up database
	psqlInfo := fmt.Sprintf("host=%s port =%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
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
	http.ListenAndServe(":8000", router)
}
