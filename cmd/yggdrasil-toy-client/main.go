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

func main() {
	// Handle command-line parameters
	genconf := flag.Bool("genconf", false, "print a new config to stdout")
	useconf := flag.Bool("useconf", false, "read HJSON/JSON config from stdin")
	useconffile := flag.String("useconffile", "", "read HJSON/JSON config from specified file path")
	normaliseconf := flag.Bool("normaliseconf", false, "use in combination with either -useconf or -useconffile, outputs your configuration normalised")
	confjson := flag.Bool("json", false, "print configuration from -genconf or -normaliseconf as JSON instead of HJSON")
	targetAddr := flag.String("target", "", "Yggdrasil address to contact")
	flag.Parse()

	var cfg *config.NodeConfig
	var state *config.NodeState
	var logger *log.Logger
	var err error

	switch {
	case *targetAddr == "":
		fmt.Println("Target flag is required.")
		return
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
	logger.Println("My node ID is:", node.Core.NodeID())
	logger.Println("My public key is:", node.Core.EncryptionPublicKey())
	logger.Println("My coords are:", node.Core.Coords())
	logger.Println("Local address:", node.Core.Address().String())
	logger.Println("Target address:", *targetAddr)

	// TODO: cleanup below this comment!
	co, err := coap.DialYggdrasil(node, *targetAddr)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	path := "/a"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := co.GetWithContext(ctx, path)

	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	log.Printf("Response payload: %v", resp.Payload())
}
