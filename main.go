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
var router = mux.NewRouter()

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
	Routes()
	http.ListenAndServe(":8000", router)
}
