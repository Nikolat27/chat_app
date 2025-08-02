package web_server

import (
	"net/http"
	"os"
	"strings"
)

func CheckCorsOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigin := getAllowedCors()

		if strings.Contains(allowedOrigin, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true") // ✅ Allow cookies
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getAllowedCors() string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		return "http://localhost:5000"
	}

	return origins
}
