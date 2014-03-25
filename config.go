package main

import (
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	port  string
	token string

	sources []Source
}

func getConfig() Config {
	var config Config

	// Read listening web port
	if len(os.Args) >= 2 {
		config.port = os.Args[1]
	} else {
		log.Println("Using default 8080 port")
		config.port = "8080"
	}

	// Read token file
	body, err := ioutil.ReadFile("token")
	if err != nil {
		log.Println("Could not read token. File 'token' should contain the diffbot token.")
		log.Fatal(err)
	}
	config.token = string(body)

	return config
}
