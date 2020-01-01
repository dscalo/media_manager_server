package main

import (
	"log"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`PONG!`))
}

func main() {
	http.HandleFunc("/ping", pingHandler)
	log.Fatal(http.ListenAndServe(":9011", nil))
}
