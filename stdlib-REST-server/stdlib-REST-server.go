package main

import (
	"log"
	"net/http"
	"time"

	"github.com/shien/restserver/stdlib-REST-server/taskserver"
)

func main() {
	mux := http.NewServeMux()
	server := taskserver.NewTaskServer()

	mux.HandleFunc("/task/", server.TaskHandler)
	mux.HandleFunc("/tag/", server.TagHandler)
	mux.HandleFunc("/due/", server.DueHandler)

	tags := []string{"BBBB", "BBBB"}
	server.Datastore.CreateTask("AAAAAAA", tags, time.Now())

	const PORT = "9090"

	log.Println("REST Server starting to listen on " + "localhost:" + PORT)
	log.Fatal(http.ListenAndServe("localhost:"+PORT, mux))

	// log.Fatal(http.ListenAndServeTLS("localhost:"+PORT, "cert.pem", "key.pem", mux))
}
