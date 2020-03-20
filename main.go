package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "gryphticon"
	dbname   = "todo"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

// Health checks the health of the API
func Health(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")                       // Logs to console API Health
	w.Header().Set("Content-Type", "application/json") // Sets the content type to JSON
	io.WriteString(w, `{"alive": true}`)               // JSON response to client
}

func main() {
	log.Info("Starting API Server")

	// Setting up database
	psqlInfo := fmt.Sprintf("host=%s port =%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Setting up the router
	router := mux.NewRouter()
	router.HandleFunc("/ping", Health).Methods("GET")
	http.ListenAndServe(":8000", router)
}
