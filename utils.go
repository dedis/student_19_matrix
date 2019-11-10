package main

import (
	"github.com/ugorji/go/codec"
	"time"
	"math/rand"
	"bytes"
	"log"
)

var (
	s1 = rand.NewSource(time.Now().UnixNano())
	r1 = rand.New(s1)
)

// Encode takes an arbitrary golang object and encodes it to JSON
func encodeJSON(val interface{}) []byte {
	var b bytes.Buffer

	var jsonH codec.Handle = new(codec.JsonHandle)
	enc := codec.NewEncoder(&b, jsonH)
	err := enc.Encode(val)
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

// Decode takes a JSON byte array and produces a golang object
func decodeJSON(pl []byte) interface{} {
	var val interface{}

	var jsonH codec.Handle = new(codec.JsonHandle)
	dec := codec.NewDecoderBytes(pl, jsonH)
	err := dec.Decode(&val)
	if err != nil {
		panic(err)
	}

	return val
}

func randSlice(n int) []byte {
	token := make([]byte, n)
	r1.Read(token)
	return token
}

func logDebug(v ...interface{}) {
	debug := true // FIXME: read from CLI arguments
	if (debug) {
		log.Println("[Debug]", v)
	}
}

func logError(v ...interface{}) {
	log.Println("Error]", v)
}
