SERVER=yggdrasil-toy-server
CLIENT=yggdrasil-toy-client
PROXY=yggdrasil-http-proxy

all: build

build: server client proxy

getDeps:
	(cd cmd/$(CLIENT); go get -v)
	(cd cmd/$(SERVER); go get -v)
	(cd cmd/$(PROXY); go get -v)

server:
	go build -o $(SERVER) cmd/$(SERVER)/main.go

client:
	go build -o $(CLIENT) cmd/$(CLIENT)/main.go

proxy:
	go build -o $(PROXY) cmd/$(PROXY)/*.go

clean:
	rm -f $(SERVER) $(CLIENT)
