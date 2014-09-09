package main

import "net/http"

func addResponseHeaders(res http.ResponseWriter, headers map[string][]string) {
	for key, value := range headers {
		res.Header().Set(key, value[0])
	}
}

func writeResponse(res http.ResponseWriter, headers map[string][]string, status int, body string) {
	addResponseHeaders(res, headers)
	res.WriteHeader(status)
	res.Write([]byte(body))
}
