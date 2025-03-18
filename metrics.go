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
