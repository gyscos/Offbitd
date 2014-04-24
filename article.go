package main

import (
	diffbot "github.com/diffbot/diffbot-go-client"
	_ "log"
	"time"
)

type Article struct {
	Title   string
	Html    string
	Url     string
	Date    string
	Author  string
	PubDate time.Time
	Read    bool
}

func wrapArticle(rawArticle *diffbot.Article) *Article {
	// log.Println("Article html: " + rawArticle.Html)
	return &Article{
		Title:   rawArticle.Title,
		Html:    rawArticle.Html,
		Date:    rawArticle.Date,
		Url:     rawArticle.Url,
		Author:  rawArticle.Author,
		PubDate: time.Now(),
		Read:    false}
}
