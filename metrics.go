package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
	DB             *DB
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {

	// HTML template with a placeholder for the variable
	htmlTemplate := `
      <html>
      <body>
          <h1>Welcome, Chirpy Admin</h1>
          <p>Chirpy has been visited %d times!</p>
      </body>
      </html>
  `

	// Format the HTML response with the variable
	response := fmt.Sprintf(htmlTemplate, cfg.fileserverHits)

	// Set the Content-Type header to "text/html; charset=utf-8"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Write the HTML response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
