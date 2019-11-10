package main

import (
	"os"
	"errors"
	"flag"
	"net/http"
	"runtime/debug"
	"sync"
	"net"
	"time"
	"fmt"

	"github.com/gologme/log"
	"github.com/yggdrasil-network/yggdrasil-go/src/config"
	"github.com/yggdrasil-network/yggdrasil-go/src/yggdrasil"

	coap "github.com/Fnux/go-coap"
	coapNet "github.com/Fnux/go-coap/net"
	toyNodes "git.sr.ht/~fnux/yggdrasil-toy-nodes"
)

var (
	// Generic variables
	err error
	logger *log.Logger
	conns = make(map[string]*net.Conn)

	// CLI Arguments
	onlyCoAP     = flag.Bool("only-coap", false, "Only proxy CoAP requests to HTTP and not the other way around")
	onlyHTTP     = flag.Bool("only-http", false, "Only proxy HTTP requests to CoAP and not the other way around")
	coapTarget   = flag.String("coap-target", "", "Force the host+port of the CoAP server to talk to")
	httpTarget   = flag.String("http-target", "http://127.0.0.1:8008", "Force the host+port of the HTTP server to talk to")
	coapPort     = flag.String("coap-port", "5683", "The CoAP port to listen on")
	coapBindHost = flag.String("coap-bind-host", "0.0.0.0", "The COAP host to listen on")
	httpBindHost = flag.String("http-bind-host", "0.0.0.0", "The HTTP host to listen on")
	httpPort     = flag.String("http-port", "8888", "The HTTP port to listen on")
	useconf      = flag.Bool("useconf", false, "read HJSON/JSON config from stdin")
	useconffile  = flag.String("useconffile", "", "read HJSON/JSON config from specified file path")
)

func httpRecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				log.Printf("Recovered from panic: %v", err)
				log.Println("Stacktrace:\n" + string(debug.Stack()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func coapRecoverWrap(h coap.Handler) coap.Handler {
	return coap.HandlerFunc(func(w coap.ResponseWriter, r *coap.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				log.Printf("Recovered from panic: %v", err)
				log.Println("Stacktrace:\n" + string(debug.Stack()))
			}
		}()
		h.ServeCOAP(w, r, )
	})
}

func main() {
	// Parse CLI arguments;
	flag.Parse()

	// Initialize logger.
	logger = log.New(os.Stdout, "", log.Flags())

	// Initialize yggdrasil-specific environment.
	var cfg *config.NodeConfig
	var state *config.NodeState

	// Process yggdrasil-related command-line arguments.
	switch {
	case *useconffile == "":
		fmt.Println("You must provide an yggdrasil configuration.")
		return
	case *useconffile != "" || *useconf:
		// Read the configuration from either stdin or from the filesystem
		normaliseconf := false
		cfg = toyNodes.ReadConfig(useconf, useconffile, &normaliseconf)
	default:
		flag.PrintDefaults()
		return
	}

	// Initialize Yggdrasil node.
	node := coapNet.YggdrasilNode{
		Core: &yggdrasil.Core{},
		Config: cfg,
	}
	state, err = node.Core.Start(node.Config, logger)
	if err != nil {
		logger.Errorln("An error occurred during Yggdrasil node startup.")
		panic(err)
	}

	// Ignore state
	_ = state

	for len(node.Core.GetPeers()) < 2 {
		logger.Println("Waiting for at least one active external peering...")
		time.Sleep(100 * time.Millisecond)
	}

	for _, peer := range node.Core.GetPeers() {
		logger.Println("Active peering:", peer.Endpoint)
	}

	// Log some basic informations.
	logger.Println("My node ID is:", node.Core.NodeID())
	logger.Println("My public key is:", node.Core.EncryptionPublicKey())
	logger.Println("My coords are:", node.Core.Coords())
	logger.Println("Local address:", node.Core.Address().String())

	// Create a wait group to keep main routine alive while HTTP and CoAP servers
	// run in separate routines.
	wg := sync.WaitGroup{}
	var h *handler

	// Start CoAP listener.
	// Listens for CoAP requests and sends out HTTP.
	if !*onlyHTTP {
		wg.Add(1)
		go func() {
			defer wg.Done()
			coapAddr := *coapBindHost + ":" + *coapPort
			log.Printf("Setting up CoAP to HTTP proxy on %s", coapAddr)
			log.Println(listenAndServe(coapAddr, "udp", coapRecoverWrap(coap.HandlerFunc(ServeCOAP))))
			log.Println("CoAP to HTTP proxy exited")
		}()
	}

	// Start HTTP listener.
	// Listens for HTTP requests and sends out CoAP.
	if !*onlyCoAP {
		wg.Add(1)
		go func() {
			defer wg.Done()
			httpAddr := *httpBindHost + ":" + *httpPort
			log.Println("Setting up HTTP to CoAP proxy on %s", httpAddr)
			log.Println(http.ListenAndServe(httpAddr, httpRecoverWrap(h)))
			log.Println("HTTP to CoAP proxy exited")
		}()
	}

	wg.Wait()

	// Close all open CoAP connections on program termination.
	for _, c := range conns {
		_ = c
		if err := (*c).Close(); err != nil {
			logError(err)
		}
	}
}
