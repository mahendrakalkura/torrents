package actions

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mahendrakalkura/torrents/go/routes"
	"github.com/mahendrakalkura/torrents/go/settings"
)

// Serve ...
func Serve() {
	addr := fmt.Sprintf("%s:%s", settings.Container.Gorilla.Hostname, settings.Container.Gorilla.Port)
	timeout := 15 * time.Second
	server := &http.Server{
		Addr:              addr,
		Handler:           routes.Connection,
		IdleTimeout:       timeout,
		ReadHeaderTimeout: timeout,
		ReadTimeout:       timeout,
		WriteTimeout:      timeout,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("actions :: Serve(...)")
	}
}
