package app

import (
	"html/template"
	"log"
	"media_manager/app/models"
	"net/http"
	"os"
	"path"
)

import "media_manager/app/uploadService"

var dirname, _ = os.Getwd()
var templates = template.Must(template.ParseFiles(path.Join(dirname, "/templates/not_found.html")))

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Env struct {
	db models.DB
}

func ChainMiddleware(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h

	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}

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
func ipLimiter(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO add logic to get ip and  continue only if internal ip
		h.ServeHTTP(w, r)
	})
}

func Run() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix(path.Join(dirname, "/static/"), fs))

	db, err := models.NewDB()
	if err != nil {
		log.Panic("Unalbe to conntect to database")
	}

	endPoints := map[string]http.HandlerFunc{
		"/":       notFoundHandler,
		"/ping":   pingHandler,
		"/upload": uploadService.UploadHandler(db),
	}

	commonMiddleware := []Middleware{
		ipLimiter,
	}

	for endPoint, fn := range endPoints {
		http.HandleFunc(endPoint, ChainMiddleware(fn, commonMiddleware...))
	}

	log.Fatal(http.ListenAndServe(":9011", nil))
}
