package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gologme/log"
	"github.com/hjson/hjson-go"
	"github.com/yggdrasil-network/yggdrasil-go/src/config"

	coap "github.com/Fnux/go-coap"
	coapNet "github.com/Fnux/go-coap/net"
	toyNodes "git.sr.ht/~fnux/yggdrasil-toy-nodes"
)

func handleA(w coap.ResponseWriter, req *coap.Request) {
	log.Printf("Got message in handleA: path=%q: %#v from %v", req.Msg.Path(), req.Msg, req.Client.RemoteAddr())
	w.SetContentFormat(coap.TextPlain)
	log.Printf("Transmitting from A")
	ctx, cancel := context.WithTimeout(req.Ctx, time.Second)
	defer cancel()
	if _, err := w.WriteWithContext(ctx, []byte("hello world")); err != nil {
		log.Printf("Cannot send response: %v", err)
	}
}

func main() {
	// Handle command-line parameters
	genconf := flag.Bool("genconf", false, "print a new config to stdout")
	useconf := flag.Bool("useconf", false, "read HJSON/JSON config from stdin")
	useconffile := flag.String("useconffile", "", "read HJSON/JSON config from specified file path")
	normaliseconf := flag.Bool("normaliseconf", false, "use in combination with either -useconf or -useconffile, outputs your configuration normalised")
	confjson := flag.Bool("json", false, "print configuration from -genconf or -normaliseconf as JSON instead of HJSON")
	flag.Parse()

	var cfg *config.NodeConfig
	var state *config.NodeState
	var logger *log.Logger
	var err error

	switch {
	case *useconffile != "" || *useconf:
		// Read the configuration from either stdin or from the filesystem
		cfg = toyNodes.ReadConfig(useconf, useconffile, normaliseconf)
		// If the -normaliseconf option was specified then remarshal the above
		// configuration and print it back to stdout. This lets the user update
		// their configuration file with newly mapped names (like above) or to
		// convert from plain JSON to commented HJSON.
		if *normaliseconf {
			var bs []byte
			if *confjson {
				bs, err = json.MarshalIndent(cfg, "", "  ")
			} else {
				bs, err = hjson.Marshal(cfg)
			}
			if err != nil {
				panic(err)
			}
			fmt.Println(string(bs))
			return
		}
	case *genconf:
		// Generate a new configuration and print it to stdout.
		fmt.Println(toyNodes.DoGenconf(*confjson))
		return
	default:
		flag.PrintDefaults()
		return
	}

	// Initialize logger
	logger = log.New(os.Stdout, "", log.Flags())

	// Initialize Yggdrasil node
	node := coapNet.YggdrasilNode{
		Config: cfg,
	}
	state, err = node.Core.Start(node.Config, logger)
	if err != nil {
		logger.Errorln("An error occurred during Yggdrasil node startup.")
		panic(err)
	}

	// Ignore state
	_ = state

	// Log some basic informations.
	logger.Println("My node ID is", node.Core.NodeID())
	logger.Println("My public key is", node.Core.EncryptionPublicKey())
	logger.Println("My coords are", node.Core.Coords())
	logger.Println("Local address ", node.Core.Address().String())

	// Launch Coap Server
	mux := coap.NewServeMux()
	mux.Handle("/a", coap.HandlerFunc(handleA))
	logger.Fatal(coap.ListenAndServeYggdrasil(node, mux))
}
