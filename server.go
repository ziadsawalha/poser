package main

import (
	"encoding/json"
	"flag"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
)

type Status struct {
	Message string
	Code    int `json:",float64"`
}

type Request struct {
	URI     string
	Method  string
	Headers map[string][]string
	Body    string
}

type Response struct {
	Headers map[string][]string
	Status  Status
	Body    map[string]interface{}
}

type Scene struct {
	Request  Request
	Response Response
}

type Scenes struct {
	Version      float64
	Interactions []Scene
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

func slices_match(slice1 []string, slice2 []string) bool {
	// Make sure the header values match
	for _, value := range slice1 {
		if !(contains(slice2, value)) {
			return false
		}
	}
	return true
}

func headers_match(expected map[string][]string, actual http.Header) bool {
	log.Printf("expected: %s\nactual: %s", expected, actual)

	// Check to see if everything in expected exists in actual
	for key, value := range expected {
		log.Printf("Testing key %s, value %s", key, value)
		if !(slices_match(value, actual[key])) {
			return false
		}
	}
	return true
}

func queries_match(query1 string, query2 string) bool {
	parsed_query1, _ := url.ParseQuery(query1)
	parsed_query2, _ := url.ParseQuery(query2)
	return reflect.DeepEqual(parsed_query1, parsed_query2)
}

func stringify(theMap map[string]interface{}) string {
	jsonified, _ := json.Marshal(theMap)
	return string(jsonified)
}

func main() {
	// Command line arguments setup
	var scenes_file = flag.String("scenes", "scenes.json", "Path to json file defining request/response pairs.")
	flag.Parse()

	// Try to parse the scenes file
	file, _ := os.Open(*scenes_file)
	decoder := json.NewDecoder(file)
	scenes := Scenes{}
	err := decoder.Decode(&scenes)
	if err != nil {
		log.Printf("%s is not a valid json scenes file.\n", *scenes_file)
		log.Fatal(err)
	}

	// Crank up Poser
	m := martini.Classic()

	m.Any("/**", func(req *http.Request) (int, string) {
		for _, scene := range scenes.Interactions {
			sceneURL, _ := url.Parse(scene.Request.URI)
			if req.Method == scene.Request.Method && req.URL.Path == sceneURL.Path {
				log.Printf("Matched method %s and URI %s\n", req.Method, req.URL.Path)

				if queries_match(req.URL.RawQuery, sceneURL.RawQuery) && headers_match(scene.Request.Headers, req.Header) {
					log.Printf("Matched query params %s", req.URL.RawQuery)
					log.Printf("Matched headers %s", scene.Request.Headers)
					return scene.Response.Status.Code, stringify(scene.Response.Body)

				} else if reflect.DeepEqual(req.Body, scene.Request.Body) && headers_match(scene.Request.Headers, req.Header) {
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
