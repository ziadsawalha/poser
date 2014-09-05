package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-martini/martini"
	"gopkg.in/yaml.v1"
)

var version = "v0.0.1" // Poser version

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
	Version      float64
	Interactions []scene
}

func contains(theSlice []string, theValue string) bool {
	log.Printf("Slice: %s, Value: %s", theSlice, theValue)
	// We don't care if the whitespace is different.
	theValue = strings.Replace(theValue, " ", "", -1)
	for _, value := range theSlice {
		value = strings.Replace(value, " ", "", -1)
		log.Printf("Comparing %s with %s", value, theValue)
		if value == theValue {
			return true
		}
	}
	log.Printf("Header did not match expectations.")
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
	log.Printf("Expected header: %s", expected)
	log.Printf("Actual   header: %s", actual)

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
	// either/both strings are empty: no match for you!
	if query1 == "" || query2 == "" {
		return false
	}

	parsedQuery1, _ := url.ParseQuery(query1)
	parsedQuery2, _ := url.ParseQuery(query2)
	return reflect.DeepEqual(parsedQuery1, parsedQuery2)
}

func bodiesMatch(expected string, actual io.ReadCloser) bool {
	actualBody, err := ioutil.ReadAll(actual)
	if err != nil {
		log.Printf("ERROR: could not read request body: comparison failed.")
		return false
	}

	log.Printf("Expected body: %s", expected)
	log.Printf("Actual body  : %s", actualBody)

	var expectedJSON, actualJSON map[string]interface{}
	json.Unmarshal([]byte(expected), &expectedJSON)
	json.Unmarshal([]byte(actualBody), &actualJSON)
	return reflect.DeepEqual(expectedJSON, actualJSON)
}

func parseScenes(scenesFilename string, allScenes *scenes) {
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

func main() {
	// Command line arguments setup
	var scenesFilename = flag.String("scenes", "scenes.json",
		"Path to json or yaml file defining request/response pairs.")
	var port = flag.String("port", "3000",
		"Port the http server should listen on. Defaults to 3000.")

	flag.Parse()
	*port = ":" + *port
	allScenes := scenes{}
	parseScenes(*scenesFilename, &allScenes)

	// Crank up Poser
	m := martini.Classic()

	m.Any("/**", func(req *http.Request) (int, string) {
		for _, scene := range allScenes.Interactions {
			sceneURL, _ := url.Parse(scene.Request.URI)
			if req.Method == scene.Request.Method && req.URL.Path == sceneURL.Path &&
				headersMatch(scene.Request.Headers, req.Header) {

				log.Printf("Matched method %s", scene.Request.Method)
				log.Printf("Matched URI %s", sceneURL.Path)
				log.Printf("Matched headers %s", scene.Request.Headers)

				if queriesMatch(req.URL.RawQuery, sceneURL.RawQuery) {
					log.Printf("Matched query params %s\n", sceneURL.RawQuery)
					return scene.Response.Status.Code, scene.Response.Body

				} else if bodiesMatch(scene.Request.Body, req.Body) {
					log.Printf("Request body matched expected.")
					return scene.Response.Status.Code, scene.Response.Body
				}
			}
		}

		// TODO(pablo): Not using a scene's provided response header... yet.

		return 501, "ERROR: Your request did not match any scenes."
	})
	log.Printf("===>>> Poser %s listening on %s <<<===", version, *port)
	log.Fatal(http.ListenAndServe(*port, m))
}
