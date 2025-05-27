package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "go-kv-store/rest-api-service/internal"
)

func main() {
    router := mux.NewRouter()
    internal.RegisterRoutes(router)

    http.Handle("/", router)
    log.Println("Starting REST API service on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}