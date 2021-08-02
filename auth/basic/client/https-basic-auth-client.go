package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", "localhost:9090/secret", "HTTPS network address")
	certFile := flag.String("certfile", "cert.pem", "trusted CA certificate")
	user := flag.String("user", "", "username")
	pass := flag.String("pass", "", "password")
	flag.Parse()

	cert, err := os.ReadFile(*certFile)

	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("Unable to parse cert from %s.", *certFile)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            certPool,
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest(http.MethodGet, "https://"+*addr, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(*user, *pass)

	rsp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer rsp.Body.Close()

	html, err := io.ReadAll(rsp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("HTTP Response Status:", rsp.Status)
	fmt.Println("Response Body:", string(html))
}
