package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// would be a database table in a real application
// And NEVER STORE PASSWORDS IN PLAINTEXT; some kind of hash should always be used
var usersPasswords = map[string][]byte{
	"shien": []byte("$2a$12$aMfFQpGSiPiYkekov7LOsu63pZFaWzmlfm1T8lvG6JFj2Bh4SZPWS"),
	"john":  []byte("$2a$12$l398tX477zeEBP6Se0mAv.ZLR8.LZZehuDgbtw2yoQeMjIyCNCsRW"),
}

func verifyUserPassword(username string, password string) bool {
	targetPassword, hasPassword := usersPasswords[username]

	if !hasPassword {
		return false
	}

	if cmpErr := bcrypt.CompareHashAndPassword(targetPassword, []byte(password)); cmpErr == nil {
		return true
	}

	return false
}

func main() {
	addr := flag.String("addr", ":9090", "HTTPS network address")
	certFile := flag.String("cerfile", "cert.pem", "certificate PEM file")
	keyFile := flag.String("key", "key.pem", "key PEM file") // private key
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rsp http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(rsp, req)
			return
		}

		fmt.Fprintf(rsp, "Proudly served with Go and HTTPS\n")
	})

	/*
		if an unauthenticated HTTP request is made to the server,
		the server adds a special header to its response: WWW-Authenticate.
		The client can then send another request, properly authenticated,
		by adding an Authorization header.
	*/
	mux.HandleFunc("/secret", func(rsp http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()

		if ok && verifyUserPassword(username, password) {
			fmt.Fprintf(rsp, "Wellcom, You get to see the secret\n")
		} else {
			rsp.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(rsp, "Unauthorized", http.StatusUnauthorized)
		}
	})

	server := http.Server{
		Addr:    *addr,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	log.Printf("Starting server on %s", *addr)
	log.Fatal(server.ListenAndServeTLS(*certFile, *keyFile))
}
