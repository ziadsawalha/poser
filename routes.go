package main

import (
	"log"
	"net/http"
	"net/url"
)

func handleAny(res http.ResponseWriter, req *http.Request) {
	for _, scene := range allScenes.Interactions {
		sceneURL, _ := url.Parse(scene.Request.URI)
		if req.Method == scene.Request.Method && req.URL.Path == sceneURL.Path &&
			headersMatch(scene.Request.Headers, req.Header) {

			log.Printf("Matched method %s", scene.Request.Method)
			log.Printf("Matched URI %s", sceneURL.Path)
			log.Printf("Matched headers %s", scene.Request.Headers)

			if queriesMatch(req.URL.RawQuery, sceneURL.RawQuery) {
				log.Printf("Matched query params %s\n", sceneURL.RawQuery)
				writeResponse(res, scene.Response.Headers, scene.Response.Status.Code,
					scene.Response.Body)
				return

			} else if bodiesMatch(scene.Request.Body, req.Body) {
				log.Printf("Request body matched expected.")
				writeResponse(res, scene.Response.Headers, scene.Response.Status.Code,
					scene.Response.Body)
				return
			}
		}
	}
	res.WriteHeader(501)
	res.Write([]byte("ERROR: request did not match any scenes."))
}
