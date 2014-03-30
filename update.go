package main

import (
	"encoding/json"
	diffbot "github.com/diffbot/diffbot-go-client"
	"log"
	"time"
)

type FrontPageItem struct {
	Title string
	URL   string
}

type FrontPage struct {
	Icon     string
	Title    string
	URL      string
	Sections []struct {
		Items   []FrontPageItem
		Primary bool
	}
}

// Replaces diffbot's broken own function
func ParseFrontpage(token string, url string, options *diffbot.Options) (*FrontPage, error) {
	body, err := diffbot.Diffbot("frontpage", token, url, options)
	if err != nil {
		log.Println("Error loading frontpage for "+url+":", err)
		return nil, err
	}

	var frontPage FrontPage
	err = json.Unmarshal(body, &frontPage)
	if err != nil {
		log.Println("Error during json parsing.")
		return nil, err
	}

	return &frontPage, nil
}

// Update a single source.
func updateSource(c *Config, source *Source, done chan struct{}) {
	// Tell him if we're done
	defer func() { done <- struct{}{} }()

	// Get the list of available articles
	frontPage, err := ParseFrontpage(c.Token, source.URL, &diffbot.Options{Fields: "*"})
	if err != nil {
		return
	}

	// Find the correct section
	for _, section := range frontPage.Sections {
		if section.Primary {
			// Add all new articles
			newItems := source.filterNewArticles(section.Items)
			for _, item := range newItems {
				article, err := diffbot.ParseArticle(c.Token, item.URL, &diffbot.Options{Fields: "*"})
				if err != nil {
					log.Println("Error!", err)
					return
				}
				log.Println("New article:", article.Title)
				source.addArticle(article)
			}
		}
	}

	log.Println("Just updated", frontPage.Title, "from", source.URL)
}

// Update the list of articles from the sources
func update(c *Config) {
	log.Println("Updating!")

	done := make(chan struct{}, len(c.Sources))

	for _, s := range c.Sources {
		go updateSource(c, s, done)
	}

	for _ = range c.Sources {
		<-done
	}

	log.Println("Update cycle complete.")
}

func updateLoop(c *Config) {
	// Periodically updates the articles from the sources
	update(c)
	for {
		select {
		case <-time.After(c.RefreshPeriod):
			update(c)
		}
	}
}
