package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-martini/martini"
)

var version = "v0.0.1"   // Poser version
var allScenes = scenes{} // All scene definitions (from scenes.go)

func main() {
	// Command line arguments setup
	var scenesFilename = flag.String("scenes", "scenes.json",
		"Path to json or yaml file defining request/response pairs.")
	var port = flag.String("port", "3000",
		"Port the http server should listen on. Defaults to 3000.")

	flag.Parse()
	parseScenes(*scenesFilename)

	// Crank up Poser
	m := martini.Classic()

	m.Any("/**", handleAny)

	log.Printf("===>>> Poser %s listening on %s <<<===", version, *port)
	log.Fatal(http.ListenAndServe(":"+*port, m))
}
