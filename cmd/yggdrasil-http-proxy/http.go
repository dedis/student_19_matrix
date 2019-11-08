package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Client for outbound HTTP requests to homeservers
var httpClient = &http.Client{}

// ???
var fedAuthPrefix = "X-Matrix origin="
var fedAuthSuffix = ",key=\"\",sig=\"\""

// handler is a struct that acts as an http.Handler, where its ServeHTTP method
// is used to handle HTTP requests
type handler struct{}

// ServeHTTP is a function implemented on handler which handles HTTP requests.
// It:
//   * Takes in an HTTP request
//   * Creates a CoAP request with carried over and headers, path, body etc.
//   * Sends the CoAP request to another proxy, retrieves the response
//   * Convert the response back into normal HTTP
//   * Returns it to the original sender
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		logger.Println("Got preflight request")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type,authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
		return
	}

	path := r.URL.Path
	logger.Println("HTTP: Got request on path %s", r.URL.Path)

	body, err := ioutil.ReadAll(r.Body)

	// Unmarshal request body JSON
	var decodedBody interface{}
	if len(body) > 0 {
		contentType := r.Header.Get("content-type")
		logger.Println("Got request with content type: %s", contentType)

		if contentType != "application/json" {
			logger.Println("Got non-json request, ignoring")

			w.WriteHeader(502)
			return
		}

		decodedBody = decodeJSON(body)
	}

	// Add authentication header to query parameters of CoAP request
	var origin *string
	if authHeader := r.Header.Get("Authorization"); len(authHeader) > 0 {
		var k, v string

		k = "access_token"
		v = strings.Replace(authHeader, "Bearer ", "", 1)

		var sep string
		if strings.Contains(path, "?") {
			sep = "&"
		} else {
			sep = "?"
		}

		path = path + sep + k + "=" + v
	}

	if len(path) == 0 {
		path = "/"
	}

	logger.Println("Final path: %s", path)

	// Send the CoAP request to another instance of the CoAP proxy and receive a response
	ctx := context.Background()
	method := "GET" // FIXME
	pl, statusCode, err := sendCoAPRequest(ctx, method, r.Host, path, decodedBody, origin)
	if err != nil {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type,authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")

	if len(pl) > 0 {
		pl = encodeJSON(pl)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(int(statusCoAPToHTTP(statusCode)))
		if _, err = w.Write(pl); err != nil {
			log.Printf("Failed to write HTTP response: %s", err.Error())
		}
	} else {
		w.WriteHeader(int(statusCoAPToHTTP(statusCode)))
	}

	logger.Println("CoAP server responded with code %s", statusCode.String())
	logger.Println("HTTP: Sending response")
}

// sendHTTPRequest is a function that sends an HTTP request to a homeserver
// either from a client or another homeserver in the case of federation.
func sendHTTPRequest(
	ctx context.Context, method string, path string, payload []byte, origin string,
) (resBody []byte, statusCode int, err error) {

	// Create the request
	url := fmt.Sprintf("%s%s", *httpTarget, path)
	logger.Println(url)
	hReq, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewReader(payload))
	if err != nil {
		logger.Println("Err preparing HTTP request", err)
		return
	}

	// Set headers
	hReq.Header.Add("Content-Type", "application/json")
	if len(origin) > 0 {
		hReq.Header.Add("Authorization", fedAuthPrefix+origin+fedAuthSuffix)
	}

	// Perform the request and receive the response
	hRes, err := httpClient.Do(hReq)
	if err != nil {
		logger.Println("Err performing HTTP request", err)
		return
	}

	// Receive the response body
	resBody, err = ioutil.ReadAll(hRes.Body)
	if err != nil {
		logger.Println("Err extracting HTTP body", err)
		return
	}

	if len(resBody) > 0 {
		contentType := hRes.Header.Get("content-type")
		logger.Println("Got response with content type: %s", contentType)

		if contentType != "application/json" {
			logger.Println("Got non-json request, ignoring")

			statusCode = 502
			resBody = []byte{}
			return
		}
	}

	statusCode = hRes.StatusCode

	return
}
