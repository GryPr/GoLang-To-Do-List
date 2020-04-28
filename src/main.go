package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"

	// Custom packages
	items "src/items"
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
	db.Debug().AutoMigrate(&items.ToDoItem{})

	log.Info("Starting API Server")

	// Setting up the router
	router := mux.NewRouter()
	router.HandleFunc("/ping", items.health).Methods("GET")
	router.HandleFunc("/todo", items.createItem).Methods("POST")
	router.HandleFunc("/todo-complete", items.getCompleteItems).Methods("GET")
	router.HandleFunc("/todo-incomplete", items.getIncompleteItems).Methods("GET")
	router.HandleFunc("/todo/{id}", items.updateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", items.deleteItem).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}
