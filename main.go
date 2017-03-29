package main

import (
	"flag"

	"github.com/mahendrakalkura/torrents/go/actions"
)

func main() {
	action := flag.String("action", "", "")
	flag.Parse()
	if *action == "query" {
		actions.Query()
	}
	if *action == "serve" {
		actions.Serve()
	}
}
