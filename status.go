package main

import (
	"net/http"
	"github.com/Fnux/go-coap"
)

// statusCoAPToHTTP is a function that converts a CoAP status code to its
// equivalent HTTP status code.
func statusCoAPToHTTP(coapCode coap.COAPCode) uint16 {
	switch coapCode {
	case coap.Content:
		return http.StatusOK
	case coap.Changed:
		return http.StatusFound
	case coap.BadRequest:
		return http.StatusBadRequest
	case coap.Unauthorized:
		return http.StatusUnauthorized
	case coap.BadOption:
		return http.StatusConflict
	case coap.Forbidden:
		return http.StatusForbidden
	case coap.NotFound:
		return http.StatusNotFound
	case coap.MethodNotAllowed:
		return http.StatusMethodNotAllowed
	case coap.RequestEntityTooLarge:
		return http.StatusTooManyRequests
	case coap.InternalServerError:
		return http.StatusInternalServerError
	case coap.BadGateway:
		return http.StatusBadGateway
	case coap.ServiceUnavailable:
		return http.StatusServiceUnavailable
	case coap.GatewayTimeout:
		return http.StatusGatewayTimeout
	default:
		logger.Println("Unsupported CoAP code %s", coapCode.String())
		return http.StatusInternalServerError
	}
}

// statusHTTPToCoAP is a function that converts an HTTP status code to its
// equivalent CoAP status code.
func statusHTTPToCoAP(httpCode int) coap.COAPCode {
	switch httpCode {
	case http.StatusOK:
		return coap.Content
	case http.StatusFound:
		return coap.Changed
	case http.StatusBadRequest:
		return coap.BadRequest
	case http.StatusUnauthorized:
		return coap.Unauthorized
	case http.StatusForbidden:
		return coap.Forbidden
	case http.StatusNotFound:
		return coap.NotFound
	case http.StatusTooManyRequests:
		return coap.RequestEntityTooLarge
	case http.StatusConflict:
		return coap.BadOption
	case http.StatusInternalServerError:
		return coap.InternalServerError
	default:
		logger.Printf("Unsupported HTTP code %d", httpCode)
		return coap.InternalServerError
	}
}
