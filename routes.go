package main

// Routes handles the routing of functions
func Routes() {
	router.HandleFunc("/ping", Health).Methods("GET", "OPTIONS")
	router.HandleFunc("/todo", CreateItem).Methods("POST", "OPTIONS")
	router.HandleFunc("/todo", GetAllItems).Methods("GET", "OPTIONS")
	router.HandleFunc("/todo-complete", GetCompleteItems).Methods("GET", "OPTIONS")
	router.HandleFunc("/todo-incomplete", GetIncompleteItems).Methods("GET", "OPTIONS")
	router.HandleFunc("/todo/{id}", UpdateCompletionItem).Methods("PUT", "OPTIONS")
	router.HandleFunc("/todo/{id}", DeleteItem).Methods("DELETE", "OPTIONS")
}
