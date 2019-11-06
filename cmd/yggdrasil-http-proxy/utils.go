package main

import (
	//"encoding/json"
	"github.com/ugorji/go/codec"
	"time"
	"math/rand"
	"bytes"
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

	// TODO: We allocate a buffer with a len(pl)*2 capacity to ensure our buffer
	// can contain all of the JSON data (as len(jsonData) > len(cborData)). This
	// is far from being the most optimised way to do it, and a more efficient
	// way of computing the maximum size of the buffer to should be
	// investigated.
	// jsn = jsn.Reset(make([]byte, 0, len(pl)*2))
	// return cbr.Reset(pl).Tojson(jsn).Bytes()
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

	// TODO: Same as above. For some reason it also blows up in some cases on
	// the JSON->CBOR way if the allocated buffer is smaller than the payload.
	// cbr = cbr.Reset(make([]byte, 0, len(pl)*2))
	// return jsn.Reset(pl).Tocbor(cbr)
}

func randSlice(n int) []byte {
	token := make([]byte, n)
	r1.Read(token)
	return token
}
