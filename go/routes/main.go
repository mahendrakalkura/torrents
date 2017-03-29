package routes

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mahendrakalkura/torrents/go/settings"
	"github.com/mahendrakalkura/torrents/go/views"
)

// Connection ...
var Connection *mux.Router

var upgrader websocket.Upgrader

func init() {
	Connection = mux.NewRouter()

	prefix := "/assets"
	directory := "assets/"
	dir := http.Dir(directory)
	fileServer := http.FileServer(dir)
	stripPrefix := http.StripPrefix(prefix, fileServer)
	Connection.PathPrefix(prefix).Handler(stripPrefix)

	Connection.HandleFunc("/favicon.ico", faviconIco)

	Connection.HandleFunc("/websockets/", websockets)
	Connection.HandleFunc("/404/", errors404).Methods("GET")
	Connection.HandleFunc("/500/", errors500).Methods("GET")
	Connection.HandleFunc("/", home).Methods("GET")

	Connection.NotFoundHandler = http.HandlerFunc(errors404)

	upgrader = websocket.Upgrader{}
}

func faviconIco(responseWriter http.ResponseWriter, request *http.Request) {
	http.ServeFile(responseWriter, request, "assets/icons/favicon.ico")
}

func websockets(responseWriter http.ResponseWriter, request *http.Request) {
	connection, connectionErr := upgrader.Upgrade(responseWriter, request, nil)
	if connectionErr != nil {
		log.Println(connection)
		return
	}
	defer connection.Close()
	for {
		messageType, messageBytes, messageErr := connection.ReadMessage()
		if messageErr != nil {
			log.Println(messageErr)
			return
		}
		messageString := string(messageBytes)
		if messageString == "start" {
			for indexInt := 1; indexInt <= 100; indexInt++ {
				indexString := strconv.Itoa(indexInt)
				messageBytes := []byte(indexString)
				err := connection.WriteMessage(messageType, messageBytes)
				if err != nil {
					log.Println(err)
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
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
