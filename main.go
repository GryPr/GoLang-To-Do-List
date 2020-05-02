package main

import (
	"fmt"
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

var psqlInfo string = fmt.Sprintf("host=%s port =%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
var db, err = gorm.Open("postgres", psqlInfo)

func main() {
	defer db.Close()

	// Setting up database

	if err != nil {
		panic(err)
	}
	if err == nil {
		fmt.Println("Sucessfully connected to the database")
	}

	//db.Debug().DropTableIfExists(&ToDoItem{})
	db.Debug().AutoMigrate(&ToDoItem{})

	log.Info("Starting API Server")

	// Setting up the router
	router := mux.NewRouter()
	router.HandleFunc("/ping", Health).Methods("GET")
	router.HandleFunc("/todo", CreateItem).Methods("POST")
	router.HandleFunc("/todo", GetAllItems).Methods("GET")
	router.HandleFunc("/todo-complete", GetCompleteItems).Methods("GET")
	router.HandleFunc("/todo-incomplete", GetIncompleteItems).Methods("GET")
	router.HandleFunc("/todo/{id}", UpdateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteItem).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}
