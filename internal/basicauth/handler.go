package basicauth

import (
	"context"
	"net/http"
)

// Authenticator will, on validation of credentials, return
// a key value pair to add to the connection context
type Authenticator interface {
	BasicAuth(un, pw string) (key, value interface{}, ok bool, err error)
}

// Handler returns a handler wrapper that applies the authenticator
func Handler(service Authenticator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", "Basic")
				http.Error(w, "bad or missing credentials", http.StatusUnauthorized)
				return
			}

			key, value, ok, err := service.BasicAuth(username, password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !ok {
				http.Error(w, "invalid credentials", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), key, value),
			))
		})
	}
}
