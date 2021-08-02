package middleware

import (
	"context"
	"net/http"

	"github.com/shien/restserver/auth/taskstore-auth/authdb"
)

/*
UserContextKey is the key in a request's context used to check if the request
has an authenticated user. The middleware will set the value of this key to
the username, if the user was properly authenticated with a password.
*/
const UserContextKey = "user"

// BasicAuth is middleware that verifies the request has appropriate basic auth
// set up with a user:password pair verified by authdb.
func BasicAuth(next http.Handler) http.Handler {
	wrappedFunc := func(rsp http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()

		if ok && authdb.VerifyUserPassword(username, password) {
			// make a key/value pair in a new Context, and pass it to the next goroutine
			newctx := context.WithValue(req.Context(), UserContextKey, username)
			next.ServeHTTP(rsp, req.WithContext(newctx))
		} else {
			rsp.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(rsp, "Unauthorized", http.StatusUnauthorized)
		}
	}

	return http.HandlerFunc(wrappedFunc)
}
