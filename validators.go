package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

func contains(theSlice []string, theValue string) bool {
	log.Printf("Slice: %s, Value: %s", theSlice, theValue)
	// We don't care if the whitespace is different.
	theValue = strings.Replace(theValue, " ", "", -1)
	for _, value := range theSlice {
		value = strings.Replace(value, " ", "", -1)
		log.Printf("Comparing %s with %s", value, theValue)
		if value == theValue {
			return true
		}
	}
	log.Printf("Header did not match expectations.")
	return false
}

func slicesMatch(slice1 []string, slice2 []string) bool {
	// Make sure the header values match
	for _, value := range slice1 {
		if !(contains(slice2, value)) {
			return false
		}
	}
	return true
}

func headersMatch(expected map[string][]string, actual http.Header) bool {
	log.Printf("Expected header: %s", expected)
	log.Printf("Actual   header: %s", actual)

	// Check to see if everything in expected exists in actual
	for key, value := range expected {
		log.Printf("Testing key %s, value %s", key, value)
		if !(slicesMatch(value, actual[key])) {
			return false
		}
	}
	return true
}

func queriesMatch(query1 string, query2 string) bool {
	// either/both strings are empty: no match for you!
	if query1 == "" || query2 == "" {
		return false
	}

	parsedQuery1, _ := url.ParseQuery(query1)
	parsedQuery2, _ := url.ParseQuery(query2)
	return reflect.DeepEqual(parsedQuery1, parsedQuery2)
}

func bodiesMatch(expected string, actual io.ReadCloser) bool {
	actualBody, err := ioutil.ReadAll(actual)
	if err != nil {
		log.Printf("ERROR: could not read request body: comparison failed.")
		return false
	}

	log.Printf("Expected body: %s", expected)
	log.Printf("Actual body  : %s", actualBody)

	var expectedJSON, actualJSON map[string]interface{}
	json.Unmarshal([]byte(expected), &expectedJSON)
	json.Unmarshal([]byte(actualBody), &actualJSON)
	return reflect.DeepEqual(expectedJSON, actualJSON)
}
