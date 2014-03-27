package main

import (
	"encoding/json"
	diffbot "github.com/diffbot/diffbot-go-client"
	"io/ioutil"
	"log"
	"strings"
)

var maxFileLength int = 60

func sanify(url string) string {
	result := strings.Replace(url, "/", "_", -1)
	if len(result) > maxFileLength {
		result = result[:maxFileLength]
	}
	return result
}

type Source struct {
	Title string
	URL   string

	NewMessages   bool
	RefreshPeriod int

	Articles []*diffbot.Article `json:"-"`
}

func (s *Source) getArticles() {
	files, err := ioutil.ReadDir(s.getDataDir())
	if err != nil {
		log.Println("Could not open articles for source "+s.Title+":", err)
	}

	for _, file := range files {
		log.Println(file.Name())
		body, err := ioutil.ReadFile(s.getDataDir() + "/" + file.Name())
		if err != nil {
			log.Println("Cannot read article file:", err)
			continue
		}

		var article diffbot.Article
		err = json.Unmarshal(body, &article)
		if err != nil {
			log.Println("Cannot unmarshal article:", err)
			continue
		}
		s.Articles = append(s.Articles, &article)
	}
}

func (s *Source) getDataDir() string {
	return "data/" + sanify(s.Title)
}

func (s *Source) filterNewArticles(items []FrontPageItem) []FrontPageItem {
	var result []FrontPageItem

	for _, item := range items {
		ok := true
		for _, article := range s.Articles {
			if article.Url == item.URL {
				ok = false
				break
			}
		}
		if ok {
			result = append(result, item)
		}
	}

	return result
}

func (s *Source) addArticle(article *diffbot.Article) {
	s.Articles = append(s.Articles, article)

	// Write article to file
	bytes, err := json.Marshal(article)
	if err != nil {
		log.Println("Error marshaling article:", err)
	}
	err = ioutil.WriteFile(s.getDataDir()+"/"+sanify(article.Url), bytes, 0600)
	if err != nil {
		log.Println("Error writing article file:", err)
	}
}
