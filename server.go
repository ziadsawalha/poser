package main

import (
	"log"
	"net/http"

	"github.com/go-martini/martini"
)

var version = "v0.0.1" // Poser version

func main() {
	parseConfig()

	// Configure Poser
	m := martini.Classic()
	m.Any("/**", handleAny)

	// Output version, port, and enabled modes information
	log.Printf("===>>> Poser %s listening on %s", version, port)

	if recordMode {
		log.Println("--->>> Poser's record/proxy mode is enabled for:")
		log.Printf("       %s", allScenes.BaseURL)
	}

	if playMode {
		log.Printf("--->>> Poser's playback mode is enabled")
	}

	// Crank up Poser
	log.Fatal(http.ListenAndServe(":"+port, m))
}
