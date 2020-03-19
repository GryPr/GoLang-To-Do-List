package main

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

// Health checks the health of the API
func Health(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")                       // Logs to console API Health
	w.Header().Set("Content-Type", "application/json") // Sets the content type to JSON
	io.WriteString(w, `{"alive": true}`)               // Sends JSON
}

func main() {
	log.Info("Starting API Server")
	router := mux.NewRouter()
	router.HandleFunc("/ping", Health).Methods("GET")
	http.ListenAndServe(":8000", router)

}
