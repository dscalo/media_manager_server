package uploadService

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"media_manager/app/models"
	"net/http"
	"os"
	"path"
	"time"
)

var dirname, _ = os.Getwd()

func UploadHandler(db *models.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		upload.MediaType = GetRootDir(upload.MimeType)

		upload.Path = path.Join(
			dirname,
			"/static/",
			upload.MediaType+"/"+upload.Name)

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

		collection := db.Client.Database("dansbrood").Collection("media")

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		media := bson.M{
			"name": upload.Name,
			"path": upload.Path,
			"type": upload.MediaType,
			"tags": upload.Tags,
			"time": time.Now(),
		}
		res, err := collection.InsertOne(ctx, media)
		if err != nil {
			log.Panic("error saving to database " + err.Error())
		}
		id := res.InsertedID

		log.Printf("added %s to database", id)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":}`))
	}
}
