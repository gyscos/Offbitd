package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Describes the state of the server.
type Config struct {

	// Port to listen on for web ui. Cannot be changed after start.
	Port string
	// Diffbot token for API access.
	Token string

	// Default duration between source refresh
	RefreshPeriod time.Duration

	// Channel to trigger a configuration dump on disk
	SyncNeeded chan struct{}

	// List of available sources
	Sources []*Source
}

func (c *Config) syncOnDisk() {
	for {
		<-c.SyncNeeded

		c.writeToFile()
	}
}

func (c *Config) requestSync() {
	c.SyncNeeded <- struct{}{}
}

func loadConfig() *Config {
	log.Println("Loading configuration.")

	if !fileExists("data") {
		err := os.Mkdir("data", os.ModeDir|0700)
		if err != nil {
			log.Println("Error creating root data dir:", err)
		}
	}

	var config Config

	config.RefreshPeriod = 5 * time.Minute
	config.Port = getPort()
	config.Token = getToken()
	config.SyncNeeded = make(chan struct{})

	// Read sources
	config.Sources = loadSources()
	for _, source := range config.Sources {
		source.prepare()
	}

	// Update the configuration on the disk when requested
	go config.syncOnDisk()

	return &config
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

func loadSources() []*Source {
	body, err := ioutil.ReadFile("data/sources.json")
	if err != nil {
		log.Println("Cannot read source file: ", err)
		return nil
	}

	source := struct{ Sources []*Source }{}
	err = json.Unmarshal(body, &source)
	if err != nil {
		log.Println("Error reading sources file: ", err)
		return nil
	}

	return source.Sources
}

func (c *Config) findSourceId(title string) int {

	for i, source := range c.Sources {
		if sanify(source.Title) == title {
			return i
		}
	}
	return -1
}

func (c *Config) findSource(title string) *Source {
	if title == "all" {
		source := &Source{Title: "All"}
		for _, s := range c.Sources {
			source.Articles = append(source.Articles, s.Articles...)
		}
		return source
	}
	i := c.findSourceId(title)
	if i == -1 {
		return nil
	}
	return c.Sources[i]
}

func (c *Config) writeToFile() {
	sourceConfig := struct{ Sources []*Source }{c.Sources}
	bytes, err := json.Marshal(sourceConfig)
	if err != nil {
		log.Println("Error marshalling sources: ", err)
	}

	err = ioutil.WriteFile("data/sources.json", bytes, 0600)
	if err != nil {
		log.Println("Error writing sources file: ", err)
	}

}

func (c *Config) addSource(url string) {
	source := makeSource(url)
	c.Sources = append(c.Sources, source)

	c.requestSync()
}

func (c *Config) removeSource(title string) {
	i := c.findSourceId(title)
	if i == -1 {
		log.Println("Error: count not find " + title)
		return
	}

	log.Println("Deleting source data: ", c.Sources[i].getDataDir())
	dataDir := c.Sources[i].getDataDir()
	// Don't inadvertently delete the root datadir !
	if dataDir != "/data/" {
		os.RemoveAll(c.Sources[i].getDataDir())
	}
	c.Sources = append(c.Sources[:i], c.Sources[i+1:]...)
	c.requestSync()

}
