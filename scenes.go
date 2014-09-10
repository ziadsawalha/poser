package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v1"
)

var allScenes = scenes{} // All scene definitions (from scenes.go)

type status struct {
	Message string
	Code    int `json:",float64"`
}

type request struct {
	URI     string
	Method  string
	Headers map[string][]string
	Body    string
}

type response struct {
	Headers map[string][]string
	Status  status
	Body    string
}

type scene struct {
	Request  request
	Response response
}

type scenes struct {
	Version float64
	BaseURL string `json:"base_url"`
	Scenes  []scene
}

func parseScenes(scenesFilename string) {
	file, _ := ioutil.ReadFile(scenesFilename)
	var err error

	if strings.HasSuffix(scenesFilename, ".json") {
		err = json.Unmarshal(file, &allScenes)
	} else if strings.HasSuffix(scenesFilename, ".yaml") {
		err = yaml.Unmarshal(file, &allScenes)
	} else {
		log.Printf("ERROR: %s does not end in '.json' or '.yaml'.", scenesFilename)
		log.Fatal(1)
	}

	if err != nil {
		log.Printf("ERROR: %s is not a valid scenes file.", scenesFilename)
		log.Fatal(err)
	}
}
