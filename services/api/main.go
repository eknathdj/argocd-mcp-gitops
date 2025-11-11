package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

func main() {
    http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    })
    http.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"now": time.Now().Format(time.RFC3339)})
    })
    log.Println("API listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
