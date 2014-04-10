package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type ListData struct {
	Sources     []*Source
	NewMessages bool
}

func handleList(w http.ResponseWriter, r *http.Request, c *Config) {
	//fmt.Fprintf(w, "LIST")

	newMessages := false
	for _, source := range c.Sources {
		if source.NewMessages > 0 {
			newMessages = true
			break
		}
	}

	data := ListData{c.Sources, newMessages}

	t, err := template.ParseFiles("templates/list.html")
	if err != nil {
		fmt.Fprintln(w, "Error loading template: ", err)
		return
	}

	t.Execute(w, data)
}

func handleView(w http.ResponseWriter, r *http.Request, c *Config) {
	target := r.URL.Path[len("/view/"):]
	source := c.findSource(target)
	if source == nil {
		// Could not find source
		log.Println("Error: could not find source " + target)
		return
	}

	t, err := template.ParseFiles("templates/view.html")
	if err != nil {
		fmt.Fprintln(w, "Error loading template: ", err)
		return
	}

	t.Execute(w, source)
}

func handleEdit(w http.ResponseWriter, r *http.Request, c *Config) {
	target := r.URL.Path[len("/edit/"):]
	source := c.findSource(target)
	if source == nil {
		// Could not find source
		log.Println("Error: could not find source " + target)
		return
	}

	t, err := template.ParseFiles("templates/edit.html")
	if err != nil {
		fmt.Fprintln(w, "Error loading template: ", err)
		return
	}

	t.Execute(w, source)
}

func handleOptions(w http.ResponseWriter, r *http.Request, c *Config) {

	t, err := template.ParseFiles("templates/options.html")
	if err != nil {
		fmt.Fprintln(w, "Error loading template: ", err)
		return
	}

	t.Execute(w, struct{ RefreshPeriod float64 }{c.RefreshPeriod.Minutes()})
}
