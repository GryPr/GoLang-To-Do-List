package main

import (
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
func Health(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")                       // Logs to console API Health
	w.Header().Set("Content-Type", "application/json") // Sets the content type to JSON
	io.WriteString(w, `{"alive": true}`)               // JSON response to client
}

func main() {
	log.Info("Starting API Server")

	// Setting up database
	psqlInfo := fmt.Sprintf("host=%s port =%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	if err == nil {
		fmt.Println("Sucessfully connected to the database")
	}
	defer db.Close()

	// Setting up the router
	router := mux.NewRouter()
	router.HandleFunc("/ping", Health).Methods("GET")
	http.ListenAndServe(":8000", router)
}
