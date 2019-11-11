package main

import (
	"errors"
	"strings"

	"github.com/Fnux/go-coap"
)

// ServeCOAP is a function that listens for CoAP requests and responds accordingly.
// It:
//   * Takes in a CoAP request
//   * Decompresses and CBOR decodes the payload if there is one
//   * Decompresses the request path and query parameters
//   * Creates an HTTP request with carried over and decompressed headers, path, body etc.
//   * Sends the HTTP request to an attached Homeserver, retrieves the response
//   * Compresses the response
//   * Returns it over CoAP to the requester
func ServeCOAP(w coap.ResponseWriter, req *coap.Request) {
	m := req.Msg
	pl := m.Payload()
	path := m.PathString()
	method := "GET" // FIXME

	if !m.IsConfirmable() {
		logDebug("Got unconfirmable message")
		return
	}

	logDebug("CoAP - %X: Got request on path %s", req.Msg.Token(), path)

	// Extract Query (???)
	var query string
	s := strings.Split(path, "?")
	if len(s) > 1 {
		query = s[1]
		path = s[0]
	} else {
		query = ""
	}
	_ = query

	logDebug("CoAP - %X: Sending HTTP request", req.Msg.Token())
	logDebug("COAP options %v", req.Msg.AllOptions())

	// Encode the raw payload into JSON
	body := encodeJSON(pl)

	// Send an HTTP request to a homeserver and receive a response
	origin := "" // FIXME
	pl, statusCode, err := sendHTTPRequest(method, path, body, origin)
	if err != nil {
		logError("Failed to sent HTTP request", err)
		return
	}

	logDebug("CoAP - %X: Got status %d", req.Msg.Token(), statusCode)
	logDebug("CoAP - %X: Sending response", req.Msg.Token())

	// Convert the receive HTTP status code to a CoAP one and add to response
	w.SetCode(statusHTTPToCoAP(statusCode))

	_, err = w.Write(pl)
	if err != nil {
			logError("Failed to return over CoAP", err)
			return
	}
}

// sendCoAPRequest is a function that sends a CoAP request to another instance
// of the CoAP proxy.
func sendCoAPRequest(method, host, path string, body interface{},
	origin *string,
) (payload []byte, statusCode coap.COAPCode, err error) {

	// Send to request's host unless the target has been forced
	var target string
	if len(*coapTarget) > 0 {
		target = *coapTarget
	} else {
		target = host
	}

	logDebug("Proxying request to %s", target)

	// Open connection to remote.
	c, err := coap.DialYggdrasil(node, target)
	if err != nil {
		logError("Error dialing: %v", err)
	}
	defer c.Close()

	// FIXME: If there is an existing connection, use it, otherwise provision a
	// new one.
	/*
	if c, exists = conns[target]; !exists || (c != nil && c.dead) {
		logger.Println("No usable connection to %s, initiating a new one", target)
		if c, err = resetConn(target); err != nil {
			return
		}
		// } else if time.Now().Add(-180 * time.Second).After(c.lastMsg) {
		// 	// Reset an existing connection if the latest message sent is older than
		// 	// go-coap's syncTimeout.
		// 	if c, err = resetConn(target); err != nil {
		// 		return
		// 	}
	} else if exists {
		logger.Println("Reusing existing connection to %s", target)
	}
	*/

	// Map for translating HTTP method codes to CoAP.
	methodCodes := map[string]coap.COAPCode{
		"GET":    coap.GET,
		"POST":   coap.POST,
		"PUT":    coap.PUT,
		"DELETE": coap.DELETE,
	}

	// Convert JSON to raw
	var bodyBytes []byte
	if body != nil {
		bodyBytes = encodeJSON(body)
	}

	logDebug("Sending %d bytes in payload", len(bodyBytes))

	// Create a new CoAP request
	req := c.NewMessage(coap.MessageParams{
		Type:      coap.Confirmable,
		Code:      methodCodes[strings.ToUpper(method)],
		MessageID: uint16(r1.Intn(100000)),
		Token:     randSlice(2),
		Payload:   bodyBytes,
	})

	if len(path) > 250 {
		// We can't send long paths, so lets bail out here
		err = errors.New("Path too long, Ignoring request: " + path)
		return
	}

	req.SetPathString(path)
	logDebug("HTTP: Sending CoAP request with token %X (path: %v)", req.Token(), path)

	// Send the CoAP request and receive a response
	logDebug("opts %v", req.AllOptions())
	res, err := c.Exchange(req)

	// Check for errors
	if err != nil {
		logError("Closing CoAP connection because of error: %v", err)

		/*
		if c, err = resetConn(target); err != nil {
			return
		}
		*/

		if res, err = c.Exchange(req); err != nil {
			logError("HTTP failed to exchange coap: %v", err)
			return
		}
	}

	// Receive the response payload
	rawPayload := res.Payload()
	logDebug("HTTP: Got response to CoAP request %X with %d bytes in response payload", res.Token(), len(rawPayload))
	pl := rawPayload

	// Keep track of the last successfully received message for connection timeout purposes
	//c.lastMsg = time.Now()

	return pl, res.Code(), err
}
