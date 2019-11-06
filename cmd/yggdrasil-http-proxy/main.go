package main

import (
	"os"
	"flag"
	"github.com/gologme/log"
)

func main() {
	onlyCoAP     := flag.Bool("only-coap", false, "Only proxy CoAP requests to HTTP and not the other way around")
	onlyHTTP     := flag.Bool("only-http", false, "Only proxy HTTP requests to CoAP and not the other way around")
	coapTarget   := flag.String("coap-target", "", "Force the host+port of the CoAP server to talk to")
	httpTarget   := flag.String("http-target", "http://127.0.0.1:8008", "Force the host+port of the HTTP server to talk to")
	coapPort     := flag.String("coap-port", "5683", "The CoAP port to listen on")
	coapBindHost := flag.String("coap-bind-host", "0.0.0.0", "The COAP host to listen on")
	httpPort     := flag.String("http-port", "8888", "The HTTP port to listen on")

	// Setting up logger.
	var logger *log.Logger
	logger = log.New(os.Stdout, "", log.Flags())

	logger.Println("FIXME")

	_ = onlyCoAP
	_ = onlyHTTP
	_ = coapTarget
	_ = httpTarget
	_ = coapPort
	_ = coapBindHost
	_ = httpPort
}
