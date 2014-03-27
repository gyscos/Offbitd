package main

import (
	"fmt"
	_ "html/template"
	"log"
	"net/http"
)

func dummyHandler(w http.ResponseWriter, r *http.Request, c *Config) {
	fmt.Fprintf(w, "HANDLING")
}

type ConfigHandler func(w http.ResponseWriter, r *http.Request, c *Config)

// Make a handler function including the given config
func makeHandler(h ConfigHandler, c *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r, c)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/list", http.StatusTemporaryRedirect)
}

func main() {

	config := getConfig()

	go updateLoop(config)

	// Visible index

	http.HandleFunc("/", mainHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/list", makeHandler(handleList, config))
	http.HandleFunc("/view/", makeHandler(handleView, config))
	http.HandleFunc("/edit/", makeHandler(dummyHandler, config))
	http.HandleFunc("/options", makeHandler(dummyHandler, config))

	// Machine-only API (via AJAX)
	http.HandleFunc("/api/source/list", makeHandler(dummyHandler, config))
	http.HandleFunc("/api/source/add/", makeHandler(handleApiAdd, config))
	http.HandleFunc("/api/source/remove/", makeHandler(handleApiRemove, config))
	http.HandleFunc("/api/article/list", makeHandler(dummyHandler, config))
	http.HandleFunc("/api/article/get/", makeHandler(dummyHandler, config))

	log.Println("Listening to port " + config.Port + "...")
	log.Fatal(http.ListenAndServe("localhost:"+config.Port, nil))
	log.Println("Exiting.")
}
