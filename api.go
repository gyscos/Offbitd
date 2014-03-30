package main

import (
	"fmt"
	_ "github.com/diffbot/diffbot-go-client"
	"log"
	"net/http"
)

func handleApiAdd(w http.ResponseWriter, r *http.Request, c *Config) {
	target := r.FormValue("target")
	log.Println("Adding new source:", target)

	// Analyse it!

	// frontPage, err := ParseFrontpage(c.Token, target, &diffbot.Options{Fields: "*"})
	//if err != nil {
	//		fmt.Fprintln(w, "Error!")
	//	return
	//}

	c.addSource(target)
	fmt.Fprintf(w, sanify(target))
}

func handleApiRemove(w http.ResponseWriter, r *http.Request, c *Config) {
	title := r.URL.Path[len("/api/source/remove/"):]
	c.removeSource(title)
}

func handleApiEdit(w http.ResponseWriter, r *http.Request, c *Config) {
	title := r.URL.Path[len("/api/source/edit/"):]
	source := c.findSource(title)
	if source == nil {
		// Could not find source
		log.Println("Error: could not find source " + title)
		fmt.Fprint(w, "error")
		return
	}

	log.Println("New Title:", r.FormValue("title"))

	source.rename(r.FormValue("title"))
	source.setURL(r.FormValue("url"))
	c.requestSync()

	fmt.Fprint(w, source.SaneTitle)
}