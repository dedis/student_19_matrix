NAME=yggdrasil-http-proxy

all: build

getDeps:
	go get -v

build:
	go build

clean:
	rm $(NAME)
