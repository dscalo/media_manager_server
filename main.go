package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var templates = template.Must(template.ParseFiles("templates/not_found.html"))

func fileExists(filename string) bool {
	info, err := os.Stat("./static/" + filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getMimeType(f multipart.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := f.Read(buffer)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buffer)

	return mimeType, nil
}

func isValidMimeType(mime string) bool {
	valid := true
	switch strings.TrimSpace(mime) {
	case "image/jpeg":
	case "image/gif":
	case "image/png":
	case "video/mpeg":
	case "video/ogg":
	default:
		valid = false
	}
	return valid
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`PONG ` + timeStamp()))
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
	r.ParseMultipartForm(1024 << 20)
	file, handler, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	if fileExists(handler.Filename) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		// todo return date of initial upload OR check of a flag to overwrite
		w.Write([]byte(`file previously uploaded`))
		return
	}
	// todo make mime checking its own function
	mime, err := getMimeType(file)

	if err != nil {
		log.Println("Unable to get mime type")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isValidMimeType(mime) {
		http.Error(w, "Invalid file type", http.StatusUnsupportedMediaType)
		return
	}
	// ******************************

	fmt.Printf("file name %+v\n", handler.Filename)
	fmt.Printf("file size %+v\n", handler.Size)
	fmt.Printf("file header %+v\n", handler.Header)

	f, err := os.OpenFile("./static/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
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
