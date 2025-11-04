package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) WriteNumberOfRequests(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileServerHits.Load()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprintf(w, "Hits: %v", hits)
}

func (cfg *apiConfig) ResetNumberOfRequests(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Write([]byte("OK"))
}

func main() {
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.Handle("/app/assets/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/healthz", http.HandlerFunc(readinessHandler))
	mux.Handle("/metrics", http.HandlerFunc(apiCfg.WriteNumberOfRequests))
	mux.Handle("/reset", http.HandlerFunc(apiCfg.ResetNumberOfRequests))

	server.ListenAndServe()
}
