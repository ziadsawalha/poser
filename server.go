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
	Uri     string
	URN     string
	Query   string
	Method  string
	Headers map[string]interface{}
	Body    string
}

type Response struct {
	Headers map[string]interface{}
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
			if req.Method == scene.Request.Method && req.URL.Path == scene.Request.URN {
				log.Printf("Matched method %s and URI %s\n", req.Method, req.URL.Path)

				if queries_match(req.URL.RawQuery, scene.Request.Query) {
					log.Printf("Matched query params %s", req.URL.RawQuery)
					return scene.Response.Status.Code, stringify(scene.Response.Body)

				} else if reflect.DeepEqual(req.Body, scene.Request.Body) {
					log.Printf("Matched request body %s", req.Body)
					return scene.Response.Status.Code, stringify(scene.Response.Body)
				}
			}
		}

		// TODO(pablo): I'm not validating against the request header yet.
		// TODO(pablo): I'm not using the response header. Should I be?

		return 501, "ERROR: Your request did not match any scenes."
	})
	m.Run()
}
