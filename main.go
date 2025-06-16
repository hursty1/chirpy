package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/hursty1/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
	secret string
	Polka_key string
}



func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	log.Printf("DB_URL: %s\n", dbURL)
	dev := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)

	TokenSecret := os.Getenv("TOKENSECRET")
	PolkaKey := os.Getenv("POLKA_KEY")
	if err != nil {
		log.Fatalf("New Error %s", err)
	}
	dbQueries := database.New(db)
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		db: dbQueries,
		platform: dev,
		secret: TokenSecret,
		Polka_key: PolkaKey,
	}
	server := http.Server {
		Addr: ":8080",
		Handler: mux,
	}

	root_handler := http.FileServer(http.Dir("."))
	handler := http.StripPrefix("/app", root_handler)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("POST /api/users", apiCfg.handleAddUser)	
	mux.HandleFunc("PUT /api/users", apiCfg.HandleUserUpdate)
	mux.HandleFunc("POST /api/login", apiCfg.HandleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandleRevoke)
	
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetAllChirps)
	mux.HandleFunc("POST /api/chirps", apiCfg.handleAddChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetChirpById)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.HandleDeleteChripById)
	
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.HandleChirpyRed)
	
	server.ListenAndServe()
}

