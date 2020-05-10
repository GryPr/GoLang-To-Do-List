package main

// Routes handles the routing of functions
func Routes() {
	router.HandleFunc("/ping", Health).Methods("GET")
	router.HandleFunc("/todo", CreateItem).Methods("POST")
	router.HandleFunc("/todo", GetAllItems).Methods("GET")
	router.HandleFunc("/todo-complete", GetCompleteItems).Methods("GET")
	router.HandleFunc("/todo-incomplete", GetIncompleteItems).Methods("GET")
	router.HandleFunc("/todo/{id}", UpdateCompletionItem).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteItem).Methods("DELETE")
}
