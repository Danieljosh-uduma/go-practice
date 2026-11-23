package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var apiHits = apiConfig{}

func main() {
	port := "8080"

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fileServer := http.FileServer(http.Dir("./"))

	mux.Handle("/app/", apiHits.middlewareMetricsInc(middlewareLog(fileServer)))

	mux.HandleFunc("POST /health", handlerReadiness)
	mux.HandleFunc("/metrics", HitsHandler)

	log.Println("Server starting on http://localhost:8080")
	log.Println("It will return 404 Not Found for every request.")

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// log.Println(r)
}

func HitsHandler(w http.ResponseWriter, r *http.Request) {
	hits := apiHits.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	response := struct {
		Hits int32 `json:"hits"`
	}{
		Hits: hits,
	}
	json.NewEncoder(w).Encode(response)
}
