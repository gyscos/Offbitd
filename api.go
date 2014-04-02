package main

import (
	"fmt"
	_ "github.com/diffbot/diffbot-go-client"
	"log"
	"net/http"
	"strconv"
	"time"
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

	source.rename(r.FormValue("title"))
	source.setURL(r.FormValue("url"))
	c.requestSync()

	fmt.Fprint(w, source.SaneTitle)
}

func handleApiRefresh(w http.ResponseWriter, r *http.Request, c *Config) {
	title := r.URL.Path[len("/api/source/refresh/"):]
	source := c.findSource(title)
	if source == nil {
		log.Println("Error: could not find source " + title)
		fmt.Fprint(w, "error")
		return
	}
	// Need a better update loop...
}

func handleApiGetArticle(w http.ResponseWriter, r *http.Request, c *Config) {
	title := r.URL.Path[len("/api/article/get/"):]
	source := c.findSource(title)
	if source == nil {
		log.Println("Error: could not find source " + title)
		fmt.Fprint(w, "error")
		return
	}

	article := source.getArticle(r.FormValue("url"))
	if article == nil {
		log.Println("Error: could not find article " + title + "/" + r.FormValue("url"))
		fmt.Fprint(w, "error")
		return
	}

	fmt.Fprintf(w, article.Text)
}

func handleApiReadArticle(w http.ResponseWriter, r *http.Request, c *Config) {
	title := r.URL.Path[len("/api/article/read/"):]
	source := c.findSource(title)
	if source == nil {
		log.Println("Error: could not find source " + title)
		fmt.Fprint(w, "error")
		return
	}

	source.markArticleRead(r.FormValue("url"))
}

func handleApiListArticles(w http.ResponseWriter, r *http.Request, c *Config) {
}

func handleApiOptions(w http.ResponseWriter, r *http.Request, c *Config) {
	p, err := strconv.ParseFloat(r.FormValue("refreshPeriod"), 64)
	if err != nil {
		log.Println("Error during period conversion!")
		return
	}

	c.RefreshPeriod = time.Duration(p * float64(time.Minute))
	log.Println("Refresh Period is now", c.RefreshPeriod)
}
