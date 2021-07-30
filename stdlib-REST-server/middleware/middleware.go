package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// Loggin middleware
func Loggin(next http.Handler) http.Handler {
	wrappedFunc := func(rsp http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(rsp, req)
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	}

	return http.HandlerFunc(wrappedFunc)
}

func PanicRecover(next http.Handler) http.Handler {
	wrappedFunc := func(rsp http.ResponseWriter, req *http.Request) {
		defer func() { // panic() called by the inner handler
			if err := recover(); err != nil {
				http.Error(rsp, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Println(string(debug.Stack()))
			}
		}()

		next.ServeHTTP(rsp, req)
	}

	return http.HandlerFunc(wrappedFunc)
}
