// offbitd project main.go
package main

import (
	"fmt"
	_ "github.com/diffbot/diffbot-go-client"
	_ "html/template"
	"log"
	"net/http"
)

func dummyHandler(w http.ResponseWriter, r *http.Request, c Config) {
	fmt.Fprintf(w, "HANDLING")
}

type ConfigHandler func(w http.ResponseWriter, r *http.Request, c Config)

func makeHandler(h ConfigHandler, c Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r, c)
	}
}

func main() {

	config := getConfig()

	go updateLoop()

	// Visible index
	http.HandleFunc("/", makeHandler(dummyHandler, config))
	http.HandleFunc("/options", makeHandler(dummyHandler, config))

	// Machine-only API (via AJAX)
	http.HandleFunc("/api/source/list", makeHandler(dummyHandler, config))
	http.HandleFunc("/api/source/add", makeHandler(dummyHandler, config))
	http.HandleFunc("/api/source/remove", makeHandler(dummyHandler, config))
	http.HandleFunc("/api/article/list", makeHandler(dummyHandler, config))
	http.HandleFunc("/api/article/get", makeHandler(dummyHandler, config))

	log.Println("Listening to port " + config.port + "...")
	log.Fatal(http.ListenAndServe("localhost:"+config.port, nil))
	log.Println("Exiting.")
}
