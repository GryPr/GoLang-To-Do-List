package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"

	// Custom packages
	items "./items"
	users "./users"
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
	items.InitDB(psqlInfo)
	users.InitDB(psqlInfo)

	// Setting up database

	if err != nil {
		panic(err)
	}
	if err == nil {
		fmt.Println("Sucessfully connected to the database")
	}

	//db.Debug().DropTableIfExists(&ToDoItem{})
	db.Debug().AutoMigrate(&items.ToDoItem{})

	log.Info("Starting API Server")

	// Setting up the router
	router := mux.NewRouter()
	router.HandleFunc("/ping", items.Health).Methods("GET")
	router.HandleFunc("/todo", items.CreateItem).Methods("POST")
	router.HandleFunc("/todo", items.GetAllItems).Methods("GET")
	router.HandleFunc("/todo-complete", items.GetCompleteItems).Methods("GET")
	router.HandleFunc("/todo-incomplete", items.GetIncompleteItems).Methods("GET")
	router.HandleFunc("/todo/{id}", items.UpdateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", items.DeleteItem).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}
