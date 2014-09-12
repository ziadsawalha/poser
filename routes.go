package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v1"
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
			URI:     allScenes.BaseURL + req.URL.RequestURI(),
			Method:  req.Method,
			Headers: req.Header,
			Body:    string(reqBody[:]),
		},
		Response: response{
			Headers: proxyRes.Header,
			Status: status{
				Message: proxyRes.Status,
				Code:    proxyRes.StatusCode,
			},
			Body: string(respBody[:]),
		},
	}

	addScene(newScene)
	var encErr error
	var encoded []byte
	if strings.HasSuffix(scenesFile, ".json") {
		encoded, encErr = json.MarshalIndent(allScenes, "", "  ")
	} else {
		encoded, encErr = yaml.Marshal(allScenes)
	}
	if encErr != nil {
		log.Printf("JSON Encoding failed: %s", encErr.Error())
	}

	file, openErr := os.OpenFile(scenesFile, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0666)
	if openErr != nil {
		log.Printf("Open file failed: %s", openErr.Error())
	}
	_, writeErr := io.WriteString(file, string(encoded))
	if writeErr != nil {
		log.Printf("Append to file failed: %s", writeErr.Error())
	}

	file.Close()
}
