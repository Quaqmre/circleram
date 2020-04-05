package auth

import (
	"net/http"
)

// Middleware is handle all endpoins with password ,  check token and pass inside
func Middleware(s Service) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ua := r.Header.Get("token")
			authrze, err := s.ParseToken(ua)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !authrze {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
