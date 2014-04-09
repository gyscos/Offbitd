package main

import (
	diffbot "github.com/diffbot/diffbot-go-client"
	"time"
)

type Article struct {
	Title   string
	Text    string
	Url     string
	Date    string
	Author  string
	PubDate time.Time
	Read    bool
}

func wrapArticle(rawArticle *diffbot.Article) *Article {
	return &Article{
		Title:   rawArticle.Title,
		Text:    rawArticle.Text,
		Date:    rawArticle.Date,
		Url:     rawArticle.Url,
		Author:  rawArticle.Author,
		PubDate: time.Now(),
		Read:    false}
}
