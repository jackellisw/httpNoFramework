package main

import (
	"net/http"
	"sync/atomic" // ensures safety for concurrency
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.Handle("GET /api/healthz", http.HandlerFunc(readinessHandler))
	mux.Handle("POST /admin/reset", http.HandlerFunc(apiCfg.handlerReset))
	mux.Handle("GET /admin/metrics", http.HandlerFunc(apiCfg.handlerMetrics))
	mux.Handle("POST /api/validate_chirp", http.HandlerFunc(handlerValidate))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
