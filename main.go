package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/jackellisw/httpNoFramework.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Println("Cannot connect to db")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
	}

	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.Handle("GET /api/healthz", http.HandlerFunc(readinessHandler))
	mux.Handle("POST /admin/reset", http.HandlerFunc(apiCfg.handlerReset))
	mux.Handle("GET /admin/metrics", http.HandlerFunc(apiCfg.handlerMetrics))
	mux.Handle("POST /api/validate_chirp", http.HandlerFunc(handlerValidate))
	mux.Handle("POST /api/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.CreateUsersHandler(w, r, apiCfg.db)
	}))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
