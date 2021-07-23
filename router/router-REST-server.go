package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shien/restserver/router/taskserver"
)

// Routing rules are not hardcoded any more, just use 3rd-party router package to handle it for us
// We just need to provide the handler functions to the routings
func main() {
	router := mux.NewRouter()
	server := taskserver.NewTaskServerForRouter()

	// By tacking a Methods call onto a route, we can easily direct different methods
	// on the same path to different handlers.
	router.HandleFunc("/task/", server.GetAllTasksHandler).Methods("GET")
	router.HandleFunc("/task/", server.CreateTaskHandler).Methods("POST")
	router.HandleFunc("/task/", server.DeleteAllTasksHandler).Methods("DELETE")

	router.HandleFunc("/task/{id:[0-9]+}", server.GetTaskHandler).Methods("GET")
	router.HandleFunc("/task/{id:[0-9]+}", server.DeleteTaskHandler).Methods("DELETE")

	router.HandleFunc("/tag/{tag}", server.TagHandler).Methods("GET")

	router.HandleFunc("/due/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}", server.DueHandler).Methods("GET")

	const PORT = "9090"

	log.Fatal(http.ListenAndServe("localhost:"+PORT, router))
}
