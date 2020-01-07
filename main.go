package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

var templates = template.Must(template.ParseFiles("templates/not_found.html"))

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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}
	// set max file size
	r.ParseMultipartForm(1024 << 20)

	file, handler, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	mime, err := GetMimeType(file)
	if err != nil {
		log.Println("Unable to get mime type")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !IsValidMimeType(mime) {
		http.Error(w, "Invalid file type", http.StatusUnsupportedMediaType)
		return
	}

	dir := GetRootDir(mime)

	if _, err := os.Stat("./static/" + dir); os.IsNotExist(err) {
		os.Mkdir("./static/"+dir, 0777)
	}

	if FileExists(handler.Filename, dir) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		// todo return date of initial upload OR check for a flag to overwrite
		w.Write([]byte(`file previously uploaded`))
		return
	}

	f, err := os.OpenFile("./static/"+dir+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, file)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": "1"}`))
}

func main() {
	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/", notFoundHandler)
	http.HandleFunc("/upload", uploadHandler)

	log.Fatal(http.ListenAndServe(":9011", nil))
}
