package main

import (
		"encoding/json"
		"flag"
		"fmt"
		"github.com/go-martini/martini"
		"log"
		"os"
)

type Response struct {
	Headers		[]map[string]string
	Status		int
	Body		string
}

type Route struct {
	Route		string
	Verb		string
	Response	Response
}

type Config struct {
	Routes	[]Route
}

func main() {
		// Command line arguments setup
        var config_file = flag.String("routes", "config.json", "Path to yaml file defining routes.")
		flag.Parse()

		// Try to parse the config file
		file, _ := os.Open(*config_file)
		decoder := json.NewDecoder(file)
		config := Config{}
		err := decoder.Decode(&config)
		if err != nil {
			fmt.Printf("%s is not a valid json configuration file.\n", *config_file)
			log.Fatal(err)
		}

		fmt.Println(config.Routes)
		fmt.Println(config.Routes[0].Route)
		fmt.Println(config.Routes[0].Verb)
		fmt.Println(config.Routes[0].Response.Headers)
		fmt.Println(config.Routes[0].Response.Status)
		fmt.Println(config.Routes[0].Response.Body)

		// Crank up Poser
        m := martini.Classic()
		for _, route := range config.Routes {
				// TODO(pablo): Obviously all routes won't be GETs...
				m.Get(route.Route, func() string {
						return route.Response.Body
				})
		}
        m.Run()
}
