package main

import (
	"log"
	"net/http"
)

func handleAny(res http.ResponseWriter, req *http.Request) {
	if playMode {
		log.Println("Entering Play Mode...")
		match, err := matchScene(res, req)

		if err == nil {
			writeResponse(res, match)
			return
		} else if !recordMode {
			res.WriteHeader(501)
			res.Write([]byte("ERROR: request did not match any scenes"))
			return
		}
	}

	// If we've made it this far, we're in record mode and need to proxy/record
	log.Println("Entering Record Mode...")
	var tempResponse response
	tempResponse.Body = "Coming soon..."
	tempResponse.Status.Code = 200
	headers := make(map[string][]string)
	headers["Content-Type"] = []string{"application/json"}
	tempResponse.Headers = headers
	writeResponse(res, tempResponse)
}
