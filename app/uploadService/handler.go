package uploadService

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var dirname, _ = os.Getwd()

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
	dirPath := path.Join(dirname, "/static/", dir) + "/"

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, 0777)
	}

	if FileExists(handler.Filename, dirPath) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		// todo return date of initial upload OR check for a flag to overwrite
		w.Write([]byte(`file previously uploaded`))
		return
	}

	f, err := os.OpenFile(dirPath+CleanFilename(handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
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
