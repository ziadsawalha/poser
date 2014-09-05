package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/go-martini/martini"
	"gopkg.in/yaml.v1"
)

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
	Body    map[string]interface{}
}

type scene struct {
	Request  request
	Response response
}

type scenes struct {
	Version      float64
	Interactions []scene
}

func contains(theSlice []string, theValue string) bool {
	log.Printf("Slice: %s, Value: %s", theSlice, theValue)
	for _, value := range theSlice {
		log.Printf("Comparing %s with %s", value, theValue)
		if value == theValue {
			return true
		}
	}
	log.Printf("Oops! Header did not match expectations")
	return false
}

func slicesMatch(slice1 []string, slice2 []string) bool {
	// Make sure the header values match
	for _, value := range slice1 {
		if !(contains(slice2, value)) {
			return false
		}
	}
	return true
}

func headersMatch(expected map[string][]string, actual http.Header) bool {
	log.Printf("expected: %s\nactual: %s", expected, actual)

	// Check to see if everything in expected exists in actual
	for key, value := range expected {
		log.Printf("Testing key %s, value %s", key, value)
		if !(slicesMatch(value, actual[key])) {
			return false
		}
	}
	return true
}

func queriesMatch(query1 string, query2 string) bool {
	parsedQuery1, _ := url.ParseQuery(query1)
	parsedQuery2, _ := url.ParseQuery(query2)
	return reflect.DeepEqual(parsedQuery1, parsedQuery2)
}

func stringify(theMap map[string]interface{}) string {
	jsonified, _ := json.Marshal(theMap)
	return string(jsonified)
}

func main() {
	// Command line arguments setup
	var scenesFile = flag.String("scenes", "scenes.json", "Path to json or yaml file defining request/response pairs.")
	flag.Parse()

	// Try to parse the scenes file
	allScenes := scenes{}
	if strings.HasSuffix(*scenesFile, ".json") {
		file, _ := os.Open(*scenesFile)
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&allScenes)
		if err != nil {
			log.Printf("%s is not a valid json scenes file.\n", *scenesFile)
			log.Fatal(err)
		}
	} else if strings.HasSuffix(*scenesFile, ".yaml") {
		file, _ := ioutil.ReadFile(*scenesFile)
		err := yaml.Unmarshal(file, &allScenes)
		if err != nil {
			log.Printf("%s is not a valid yaml scenes file.\n", *scenesFile)
			log.Fatal(err)
	    }
	} else {
		log.Printf("%s does not end in '.json' or '.yaml'.\n", *scenesFile)
		log.Fatal(1)
	}

	// Crank up Poser
	m := martini.Classic()

	m.Any("/**", func(req *http.Request) (int, string) {
		for _, scene := range allScenes.Interactions {
			sceneURL, _ := url.Parse(scene.Request.URI)
			if req.Method == scene.Request.Method && req.URL.Path == sceneURL.Path {
				log.Printf("Matched method %s and URI %s\n", req.Method, req.URL.Path)

				if queriesMatch(req.URL.RawQuery, sceneURL.RawQuery) &&
					headersMatch(scene.Request.Headers, req.Header) {

					log.Printf("Matched query params %s", req.URL.RawQuery)
					log.Printf("Matched headers %s", scene.Request.Headers)
					return scene.Response.Status.Code, stringify(scene.Response.Body)

				} else if reflect.DeepEqual(req.Body, scene.Request.Body) &&
					headersMatch(scene.Request.Headers, req.Header) {

					log.Printf("Matched request body %s", req.Body)
					log.Printf("Matched headers %s", scene.Request.Headers)
					return scene.Response.Status.Code, stringify(scene.Response.Body)
				}
			}
		}

		// TODO(pablo): Not using a scene's provided response header... yet.

		return 501, "ERROR: Your request did not match any scenes."
	})
	m.Run()
}
