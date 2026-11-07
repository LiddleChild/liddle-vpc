package middleware

import (
	"crypto/subtle"
	"net/http"
)

type apiKey struct {
	key string
}

func NewAPIKey(key string) apiKey {
	return apiKey{
		key: key,
	}
}

func (m apiKey) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")

		if subtle.ConstantTimeCompare([]byte(key), []byte(m.key)) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
