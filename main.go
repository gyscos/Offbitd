// offbitd project main.go
package main

import (
	"fmt"
	_ "github.com/diffbot/diffbot-go-client"
	_ "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	port  string
	token string
}

func getConfig() Config {
	var config Config

	// Read listening web port
	if len(os.Args) >= 2 {
		config.port = os.Args[1]
	} else {
		log.Println("Using default 8080 port")
		config.port = "8080"
	}

	// Read token file
	body, err := ioutil.ReadFile("token")
	if err != nil {
		log.Fatal("Error!! ", err)
	}
	config.token = string(body)

	return config
}

// Update the list of articles from the sources
func update() {

}

func updateLoop(ticks <-chan time.Time) {
	// Periodically updates the articles from the sources
	select {
	case <-ticks:
		// Update list
		update()
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HANDLING")
}

func main() {

	config := getConfig()

	ticks := make(chan time.Time)
	go updateLoop(ticks)

	// Visible index
	http.HandleFunc("/", handler)

	// Machine-only API (via AJAX)
	http.HandleFunc("/source/list", handler)
	http.HandleFunc("/source/add", handler)
	http.HandleFunc("/source/remove", handler)
	http.HandleFunc("/article/list", handler)
	http.HandleFunc("/article/get", handler)

	log.Println("Listening to port " + config.port + "...")
	log.Fatal(http.ListenAndServe("localhost:"+config.port, nil))
	log.Println("Exiting.")
}
