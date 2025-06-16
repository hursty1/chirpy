package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)


func healthzHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	io.WriteString(rw, "OK")
}

// Handler to reset hit count
func (cfg *apiConfig) resetHandler(rw http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		responseWithError(rw, 403, "Forbidden")
	}
	cfg.db.DeleteAllUsers(context.Background())
	log.Println("Deleting all users")
	cfg.fileserverHits.Store(0)
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(rw, "Hits counter reset\n")
}

func (cfg *apiConfig) metricsHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	count := cfg.fileserverHits.Load()

	html := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, count)

	io.WriteString(rw, html)
}
