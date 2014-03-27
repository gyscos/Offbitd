package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	Port  string
	Token string

	RefreshPeriod time.Duration

	Sources []Source
}

// Type used to marshal into config file
type SourceConfig struct {
	Sources []Source
}

func getPort() string {
	// Read listening web port
	if len(os.Args) >= 2 {
		return os.Args[1]
	} else {
		log.Println("Using default 8080 port")
		return "8080"
	}
}

func getToken() string {
	body, err := ioutil.ReadFile("token")
	if err != nil {
		log.Println("Could not read token. File 'token' should contain the diffbot token.")
		log.Fatal(err)
	}
	token := string(body)
	if token[len(token)-1] == '\n' {
		token = token[:len(token)-1]
	}
	return token
}

func getSources() []Source {
	body, err := ioutil.ReadFile("data/sources.json")
	if err != nil {
		log.Println("Cannot read source file: ", err)
		return nil
	}

	var source SourceConfig
	err = json.Unmarshal(body, &source)
	if err != nil {
		log.Println("Error reading sources file: ", err)
		return nil
	}

	return source.Sources
}

func getConfig() *Config {
	var config Config

	config.RefreshPeriod = 5 * time.Minute
	config.Port = getPort()
	config.Token = getToken()
	// Read sources

	config.Sources = getSources()
	for _, source := range config.Sources {
		source.getArticles()
	}

	return &config
}

func (c *Config) FindSourceId(title string) int {

	for i, source := range c.Sources {
		if source.Title == title {
			return i
		}
	}
	return -1
}

func (c *Config) FindSource(title string) *Source {
	if title == "all" {
		source := &Source{Title: "All"}
		for _, s := range c.Sources {
			source.Articles = append(source.Articles, s.Articles...)
		}
		return source
	}
	i := c.FindSourceId(title)
	if i == -1 {
		return nil
	}
	return &c.Sources[i]
}

func (c *Config) writeToFile() {
	sourceConfig := SourceConfig{c.Sources}
	bytes, err := json.Marshal(sourceConfig)
	if err != nil {
		log.Println("Error marshalling sources: ", err)
	}

	err = ioutil.WriteFile("data/sources.json", bytes, 0600)
	if err != nil {
		log.Println("Error writing sources file: ", err)
	}

}

func (c *Config) AddSource(source Source) {
	c.Sources = append(c.Sources, source)

	err := os.Mkdir(source.getDataDir(), os.ModeDir|0700)
	if err != nil {
		log.Println("Error creating source directory:", err)
	}

	c.writeToFile()
}

func (c *Config) RemoveSource(title string) {
	i := c.FindSourceId(title)
	if i == -1 {
		log.Println("Error: count not find " + title)
		return
	}
	os.RemoveAll(c.Sources[i].getDataDir())
	c.Sources = append(c.Sources[:i], c.Sources[i+1:]...)
	c.writeToFile()

}
