package uploadService

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var dirname, _ = os.Getwd()

type Upload struct {
	name      string
	path      string
	mimeType  string
	mediaType string
	tags      []string
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
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

	upload := Upload{
		name:     CleanFilename(handler.Filename),
		mimeType: GetMimeType(file),
	}

	defer file.Close()

	if upload.mimeType == "" {
		log.Println("Unable to get mime type")
		http.Error(w, "Invalid mime type", http.StatusInternalServerError)
		return
	}

	if !IsValidMimeType(upload.mimeType) {
		http.Error(w, "Invalid file type", http.StatusUnsupportedMediaType)
		return
	}

	upload.path = path.Join(dirname, "/static/", GetRootDir(upload.mimeType)) + "/" + upload.name

	if FileExists(upload.path) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		// todo return date of initial upload OR check for a flag to overwrite
		w.Write([]byte(`file previously uploaded`))
		return
	}

	f, err := os.OpenFile(upload.path, os.O_WRONLY|os.O_CREATE, 0666)
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
