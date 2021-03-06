package main

import (
	"flag"
	"log"
)

var playMode bool     // false = Play Mode disabled (default)
var port string       // HTTP port
var proxy string      // For record mode, the URL Poser should proxy to
var recordMode bool   // false = Record Mode disabled (default)
var scenesFile string // The json/yaml file containing all scenes

func parseConfig() {
	flag.StringVar(&scenesFile, "scenes", "scenes.json",
		"Path to json/yaml scenes file. Default is \"\"")
	flag.StringVar(&port, "port", "3000", "HTTP port. Defaults to 3000")
	flag.BoolVar(&playMode, "play", false, "Enable play mode. Default is `true`")
	flag.BoolVar(&recordMode, "record", false, "Enable record mode. Default is `false`")
	flag.StringVar(&proxy, "proxy", "",
		"The URL for which Poser should act as proxy. Default is \"\"")

	flag.Parse()
	parseScenes(scenesFile)

	if proxy != "" {
		// Command line takes precedence: overwrite BaseURL and Version
		allScenes.BaseURL = proxy
		allScenes.Version = 1
	}

	if recordMode && allScenes.BaseURL == "" {
		log.Fatal("A proxy is required for record mode. Provide a valid URL " +
			"either by adding it to your scenes file under the \"base_url\" " +
			"key or pass it in using the -proxy command line parameter.")
	}

	if !(recordMode || playMode) {
		log.Fatal("Please enable at least one mode and try again (-record and/or -play).")
	}
}
