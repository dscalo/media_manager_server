package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/not_found.html"))

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`PONG!`))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "not_found.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))

	http.HandleFunc("/ping", pingHandler)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", notFoundHandler)

	log.Fatal(http.ListenAndServe(":9011", nil))
}
