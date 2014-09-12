package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
		log.Fatal("Scenes file must end in either \".json\" or \".yaml\"")
	}

	if err != nil && playMode {
		log.Fatal("A valid scenes file is required (" + scenesFilename + " is not valid)")
	}
}

func matchScene(res http.ResponseWriter, req *http.Request) (match response, err error) {
	for _, scene := range allScenes.Scenes {
		sceneURL, _ := url.Parse(scene.Request.URI)
		if req.Method == scene.Request.Method && req.URL.Path == sceneURL.Path &&
			headersMatch(scene.Request.Headers, req.Header) {

			log.Printf("Matched method %s", scene.Request.Method)
			log.Printf("Matched URI %s", sceneURL.Path)
			log.Printf("Matched headers %s", scene.Request.Headers)

			if queriesMatch(req.URL.RawQuery, sceneURL.RawQuery) {
				log.Printf("Matched query params %s\n", sceneURL.RawQuery)
				return scene.Response, err

			} else if bodiesMatch(scene.Request.Body, req.Body) {
				log.Println("Request body matched expected.")
				return scene.Response, err
			}
		}
	}
	log.Println("No scene was found that matched the request")
	return match, errors.New("matchScene: no scene matched request")
}

func addScene(newScene scene) {
	n := len(allScenes.Scenes)
	if n == cap(allScenes.Scenes) {
		newScenes := make([]scene, len(allScenes.Scenes), 2*len(allScenes.Scenes)+1)
		copy(newScenes, allScenes.Scenes)
		allScenes.Scenes = newScenes
	}
	allScenes.Scenes = allScenes.Scenes[0 : n+1]
	allScenes.Scenes[n] = newScene
}
