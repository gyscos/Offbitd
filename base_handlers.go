package main

import (
	"fmt"
	_ "github.com/diffbot/diffbot-go-client"
	"html/template"
	"net/http"
)

type SourceData struct {
	Title       string
	URL         string
	NewMessages bool
}

type ListData struct {
	Sources     []SourceData
	NewMessages bool
}

func handleList(w http.ResponseWriter, r *http.Request, c Config) {
	//fmt.Fprintf(w, "LIST")

	var data ListData
	data.Sources = append(data.Sources, SourceData{"Engadget", "http://engadget.com", true})
	data.Sources = append(data.Sources, SourceData{"The Verge", "http://theverge.com", false})
	data.NewMessages = true

	t, err := template.ParseFiles("templates/list.html")
	if err != nil {
		fmt.Fprintln(w, "Error loading template: ", err)
		return
	}

	t.Execute(w, data)
}
