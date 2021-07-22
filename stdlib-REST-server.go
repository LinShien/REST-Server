package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	server := NewTaskServer()

	mux.HandleFunc("/task/", server.taskHandler)
	mux.HandleFunc("/tag/", server.tagHandler)

	tags := []string{"BBBB", "BBBB"}
	server.datastore.CreateTask("AAAAAAA", tags, time.Now())

	PORT := "9090"

	log.Println("REST Server starting to listen on " + "localhost:" + PORT)
	// log.Fatal(http.ListenAndServeTLS("localhost:"+PORT, "cert.pem", "key.pem", mux))
	log.Fatal(http.ListenAndServe("localhost:"+PORT, mux))
}
