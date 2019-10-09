package main

import (
  "os"
  "github.com/gologme/log"
  "github.com/yggdrasil-network/yggdrasil-go/src/config"
  "github.com/yggdrasil-network/yggdrasil-go/src/yggdrasil"
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

  // Listen for incoming events
  listener, err := n.core.ConnListen()
  if err != nil {
    n.log.Errorln("An error occured setting up the Yggdrasil listener.")
    panic(err)
  }

  for {
    conn, err := listener.Accept()
    if err != nil {
      n.log.Errorln("An error occured on incoming connection.")
      panic(err)
    }

    n.log.Println("New connection!")
    _ = conn
  }
}
