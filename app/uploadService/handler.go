package uploadService

import (
	"io"
	"log"
	"media_manager/app/models"
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
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	upload := models.Upload{
		Name:     CleanFilename(handler.Filename),
		MimeType: GetMimeType(file),
		Tags:     ParseTags(r.Form.Get("tags")),
	}

	if !IsValidMimeType(upload.MimeType) {
		log.Println("Unable to get mime type")
		http.Error(w, "Invalid mime type", http.StatusInternalServerError)
		return
	}

	upload.Path = path.Join(
		dirname,
		"/static/",
		GetRootDir(upload.MimeType)) + "/" + upload.Name

	if FileExists(upload.Path) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		// todo return date of initial upload OR check for a flag to overwrite
		w.Write([]byte(`file previously uploaded`))
		return
	}

	f, err := os.OpenFile(upload.Path, os.O_WRONLY|os.O_CREATE, 0666)
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
