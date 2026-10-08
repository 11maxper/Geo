package mymiddleware

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func CustomAuthenticator(tokenAuth *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}
			if claims == nil {
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
