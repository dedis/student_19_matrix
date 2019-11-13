# Yggdrasil HTTP Proxy

This application proxy HTTP requests over
[Yggdrasil](https://yggdrasil-network.github.io/) using the [CoAP
protocol](http://coap.technology/). It uses a [patched version of
go-coap](https://github.com/Fnux/go-coap) which natively support Yggdrasil
connections.

## Building

The projects depends on [Go](https://golang.org/) `>= 1.13` and is built as
follow:

```
go get -v # or `make getDeps`
go build # or `make`
```

## Usage

FIXME: provides examples and details setup.

```
./yggdrasil-http-proxy -help
Usage of ./yggdrasil-http-proxy:
  -coap-target string
        Force the host+port of the CoAP server to talk to
  -http-bind-host string
        The HTTP host to listen on (default "0.0.0.0")
  -http-port string
        The HTTP port to listen on (default "8888")
  -http-target string
        Force the host+port of the HTTP server to talk to (default "http://127.0.0.1:8008")
  -only-coap
        Only proxy CoAP requests to HTTP and not the other way around
  -only-http
        Only proxy HTTP requests to CoAP and not the other way around
  -useconf
        read HJSON/JSON config from stdin
  -useconffile string
        read HJSON/JSON config from specified file path
```

## Acknowledgment & Licensing

The `config.go` file has been imported almost as-in from
[yggdrasil-go](https://github.com/yggdrasil-network/yggdrasil-go) (LGPLv3)
while the proxy logic is heavily inspired from [matrix-org's
coap-proxy](https://github.com/matrix-org/coap-proxy) (GPL-3.0).

I haven't decided on anything for the remaining parts yet, please complain if
you need one.
