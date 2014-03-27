package main

import (
	"fmt"
	diffbot "github.com/diffbot/diffbot-go-client"
	"net/http"
)

func handleApiAdd(w http.ResponseWriter, r *http.Request, c *Config) {
	url := r.URL.Path[len("/api/source/add/"):]

	// Analyse it!

	frontPage, err := ParseFrontpage(c.Token, url, &diffbot.Options{Fields: "*"})
	if err != nil {
		fmt.Fprintln(w, "Error!")
		return
	}

	c.AddSource(Source{URL: url, Title: frontPage.Title})
	fmt.Fprintf(w, frontPage.Title)
}

func handleApiRemove(w http.ResponseWriter, r *http.Request, c *Config) {
	title := r.URL.Path[len("/api/source/remove/"):]
	c.RemoveSource(title)
}
