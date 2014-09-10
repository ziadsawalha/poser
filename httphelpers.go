package main

import "net/http"

func addResponseHeaders(res http.ResponseWriter, headers map[string][]string) {
	for key, value := range headers {
		res.Header().Set(key, value[0])
	}
}

func writeResponse(res http.ResponseWriter, response response) {
	addResponseHeaders(res, response.Headers)
	res.WriteHeader(response.Status.Code)
	res.Write([]byte(response.Body))
}
