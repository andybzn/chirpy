package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	bodyContent := fmt.Sprintf("<html>\n    <body>\n        <h1>Welcome, Chirpy Admin</h1>\n        <p>Chirpy has been visited %d times!</p>\n    </body>\n</html>", cfg.fileserverHits.Load())
	w.Write([]byte(bodyContent))
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
