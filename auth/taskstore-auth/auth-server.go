package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/shien/restserver/auth/taskstore-auth/middleware"
	"github.com/shien/restserver/auth/taskstore-auth/taskserver"
)

func main() {
	certFile := flag.String("certfile", "cert.pem", "certificate PEM file")
	keyFile := flag.String("keyfile", "key.pem", "key PEM file")
	flag.Parse()

	router := mux.NewRouter()
	router.StrictSlash(true)
	taskServer := taskserver.NewTaskServerForRouter()

	router.Handle("/task/", middleware.BasicAuth(http.HandlerFunc(taskServer.CreateTaskHandler))).Methods("POST")

	router.HandleFunc("/task/", taskServer.GetAllTasksHandler).Methods("GET")
	router.HandleFunc("/task/", taskServer.DeleteAllTasksHandler).Methods("DELETE")

	router.HandleFunc("/task/{id:[0-9]+}", taskServer.GetTaskHandler).Methods("GET")
	router.HandleFunc("/task/{id:[0-9]+}", taskServer.DeleteTaskHandler).Methods("DELETE")

	router.HandleFunc("/tag/{tag}", taskServer.TagHandler).Methods("GET")

	router.HandleFunc("/due/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}", taskServer.DueHandler).Methods("GET")

	router.Use(func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	})

	addr := "localhost:9090"
	server := &http.Server{
		Addr:    addr,
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	log.Printf("Starting server on %s", addr)
	log.Fatal(server.ListenAndServeTLS(*certFile, *keyFile))
}
