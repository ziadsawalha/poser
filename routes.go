package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

	proxyURL, _ := url.Parse(allScenes.BaseURL)
	req.RequestURI = ""
	req.URL.Scheme = proxyURL.Scheme
	req.URL.Host = proxyURL.Host
	req.Host = proxyURL.Host

	client := &http.Client{}
	proxyRes, err := client.Do(req)

	if err != nil {
		res.WriteHeader(501)
		res.Write([]byte("Proxy call failed: " + err.Error()))
	}

	addResponseHeaders(res, proxyRes.Header)
	res.WriteHeader(proxyRes.StatusCode)
	respBody, _ := ioutil.ReadAll(proxyRes.Body)
	res.Write(respBody)

	reqBody, _ := ioutil.ReadAll(req.Body)

	newScene := scene{
		Request: request{
			URI:     "",
			Method:  req.Method,
			Headers: req.Header,
			Body:    string(reqBody[:]),
		},
		Response: response{
			Headers: proxyRes.Header,
			Status: status{
				Message: "",
				Code:    proxyRes.StatusCode,
			},
			Body: string(respBody[:]),
		},
	}
	log.Println(newScene)
}
