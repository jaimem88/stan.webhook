package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// HTTP config
type HTTP struct {
	ListenPort string `json:"listen_port,omitempty"`
}

var config = struct {
	Environment string `json:"environment,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	HTTP        *HTTP  `json:"http,omitempty"`
}{
	Environment: "",
	Hostname:    "",
	HTTP: &HTTP{
		ListenPort: os.Getenv("PORT"),
	},
}

func writeDefaultConfig(location string) {
	f, err := os.Create(location)
	if err != nil {
		log.Fatalln("Couldn't open", location)
	}

	data, _ := json.MarshalIndent(config, "", "  ")
	_, err = f.Write(data)
	if err != nil {
		log.Fatalln("Couldn't write to", location)
	}
}

func loadConfig(location string) {
	raw, err := ioutil.ReadFile(location)
	if err != nil {
		log.Fatalln("Couldn't open ", location)
	}

	err = json.Unmarshal(raw, &config)
	if err != nil {
		log.Fatalln("Couldn't understand config in", location, "-", err)
	}
}
