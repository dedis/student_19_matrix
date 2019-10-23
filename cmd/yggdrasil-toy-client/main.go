package main

import (
  "os"
  "time"
  "context"
  "github.com/gologme/log"
  "github.com/yggdrasil-network/yggdrasil-go/src/config"
  "github.com/yggdrasil-network/yggdrasil-go/src/yggdrasil"

  coap "github.com/Fnux/go-coap"
  coapNet "github.com/Fnux/go-coap/net"
)

// Defines an Yggdrasil node.
type node struct {
  core   yggdrasil.Core
  config *config.NodeConfig
  state  *config.NodeState
  log    *log.Logger
}

func initLocalNode() node {
  n := node{}
  n.log = log.New(os.Stdout, "", log.Flags())
  n.config = config.GenerateConfig()

  return n
}

func main() {
  var err error

  // Initialize local Yggdrasil node.
  n := initLocalNode()

  // Start node.
  n.log.Println("Starting Yggdrasil node.")
  n.state, err = n.core.Start(n.config, n.log)
  if err != nil {
    n.log.Errorln("An error occurred during startup")
    panic(err)
  }

  // Log some basic informations.
  n.log.Println("My node ID is", n.core.NodeID())
  n.log.Println("My public key is", n.core.EncryptionPublicKey())
  n.log.Println("My coords are", n.core.Coords())

  // Connect to the global Yggdrasil network.
  n.log.Println("Connecting to global Network.")
  // -- From https://github.com/yggdrasil-network/public-peers
  swissBayPeer := "tcp://77.56.134.244:34962"
  n.core.AddPeer(swissBayPeer, "")

  n.log.Println("Local address ", n.core.Address().String())

  // TODO: Send HTTP query to toy server
  target := "303:60d4:3d32:a2b9::4" // Some kind of yggdrasil-enabled forum

	yggdrasilCoapNode := coapNet.YggdrasilNode{ Core: n.core, Config: n.config }
	co, err := coap.DialYggdrasil(yggdrasilCoapNode, target)
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
