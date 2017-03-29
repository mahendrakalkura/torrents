package views

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

// Assets ...
var Assets map[string]int

// Templates ...
var Templates map[string]*template.Template

func init() {
	initAssets()
	initTemplates()
}

func initAssets() {
	Assets = make(map[string]int)
	files := []string{
		"assets/compressed.css",
		"assets/compressed.min.css",
		"assets/compressed.js",
		"assets/compressed.min.js",
	}
	for _, file := range files {
		stat, statErr := os.Stat(file)
		if statErr != nil {
			log.Fatalln("views :: initAssets(...)")
		}
		modTime := stat.ModTime()
		secondsInt64 := modTime.Unix()
		secondsInt := int(secondsInt64)
		Assets[file] = secondsInt
	}
}

func initTemplates() {
	funcs := template.FuncMap{
		"assets": assets,
	}

	layouts, layoutsErr := filepath.Glob("resources/html/*.html")
	if layoutsErr != nil {
		log.Fatalln("templates :: initTemplates(...) - filepath.Glob(...) - #1")
	}

	views, viewsErr := filepath.Glob("resources/html/routes/*.html")
	if viewsErr != nil {
		log.Fatalln("templates :: initTemplates(...) - filepath.Glob(...) - #1")
	}

	Templates = make(map[string]*template.Template)
	for _, view := range views {
		files := append(layouts, view)
		parseFiles, parseFilesErr := template.New("layout").Funcs(funcs).ParseFiles(files...)
		if parseFilesErr != nil {
			log.Println(parseFilesErr)
			log.Fatalln("templates :: initTemplates(...) - templates.New(...).Funcs(...).ParseFiles(...)")
		}
		must := template.Must(parseFiles, parseFilesErr)
		Templates[view] = must
	}
}

func assets(file string) string {
	hash, ok := Assets[file]
	if !ok {
		log.Fatalln("templates :: assets(...)")
	}
	return fmt.Sprintf("%s?timestamp=%d", file, hash)
}
