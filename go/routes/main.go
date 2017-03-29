package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mahendrakalkura/torrents/go/settings"
	"github.com/mahendrakalkura/torrents/go/views"
)

// Connection ...
var Connection *mux.Router

func init() {
	Connection = mux.NewRouter()

	prefix := "/assets"
	directory := "assets/"
	dir := http.Dir(directory)
	fileServer := http.FileServer(dir)
	stripPrefix := http.StripPrefix(prefix, fileServer)
	Connection.PathPrefix(prefix).Handler(stripPrefix)

	Connection.HandleFunc("/favicon.ico", faviconIco)

	Connection.HandleFunc("/404/", errors404).Methods("GET")
	Connection.HandleFunc("/500/", errors500).Methods("GET")
	Connection.HandleFunc("/", home).Methods("GET")
	Connection.HandleFunc("/", xhr).Methods("POST")

	Connection.NotFoundHandler = http.HandlerFunc(errors404)
}

func faviconIco(responseWriter http.ResponseWriter, request *http.Request) {
	http.ServeFile(responseWriter, request, "assets/icons/favicon.ico")
}

func errors404(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusNotFound)
	data := map[string]interface{}{
		"settings": settings.Container,
	}
	err := views.Templates["resources/html/routes/404.html"].Execute(responseWriter, data)
	if err != nil {
		log.Println(err)
	}
}

func errors500(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusInternalServerError)
	data := map[string]interface{}{
		"settings": settings.Container,
	}
	err := views.Templates["resources/html/routes/500.html"].Execute(responseWriter, data)
	if err != nil {
		log.Println(err)
	}
}

func home(responseWriter http.ResponseWriter, request *http.Request) {
	data := map[string]interface{}{
		"settings": settings.Container,
	}
	err := views.Templates["resources/html/routes/home.html"].Execute(responseWriter, data)
	if err != nil {
		log.Println(err)
	}
}

func xhr(responseWriter http.ResponseWriter, request *http.Request) {
}
