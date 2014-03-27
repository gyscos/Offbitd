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

// Update the list of articles from the sources
func update(c *Config) {
	log.Println("Updating!")
	opt := &diffbot.Options{Fields: "*"}

	done := make(chan struct{}, len(c.Sources))

	for _, s := range c.Sources {
		source := s
		go func() {
			defer func() { done <- struct{}{} }()
			frontPage, err := ParseFrontpage(c.Token, source.URL, opt)
			if err != nil {
				return
			}

			for _, section := range frontPage.Sections {
				if section.Primary {
					newItems := s.filterNewArticles(section.Items)
					for _, item := range newItems {
						log.Println(item.URL)
						article, err := diffbot.ParseArticle(c.Token, item.URL, opt)
						if err != nil {
							log.Println("Error!", err)
							return
						}
						source.addArticle(article)
					}
				}
			}

			log.Println("Title:", frontPage.Title, "from", source.URL)
		}()
	}
	for _ = range c.Sources {
		<-done
	}
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
