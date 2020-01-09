package app

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

import "media_manager/app/uploadService"

var dirname, _ = os.Getwd()

var templates = template.Must(template.ParseFiles(path.Join(dirname, "/templates/not_found.html")))

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`PONG`))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "not_found.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Run() {
	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix(path.Join(dirname, "/static/"), fs))
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/", notFoundHandler)
	http.HandleFunc("/upload", uploadService.UploadHandler)

	log.Fatal(http.ListenAndServe(":9011", nil))
}
