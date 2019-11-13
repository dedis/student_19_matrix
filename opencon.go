package main

import (
	"github.com/Fnux/go-coap"
	"fmt"
)

var conns = make(map[string]*coap.ClientConn)

// FIXME: should this be handled on the go-coap-side?
// FIXME: cleanup broken connections
func dial(target string) (c *coap.ClientConn, err error) {
	if conns[target] != nil {
		return conns[target], err
	}

	c, err = coap.DialYggdrasil(node, target)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}

	conns[target] = c
	return c, nil
}
