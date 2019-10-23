SERVER=yggdrasil-toy-server
CLIENT=yggdrasil-toy-client

all: build

build: server client

getDeps:
	(cd cmd/yggdrasil-toy-server; go get -v)
	(cd cmd/yggdrasil-toy-client; go get -v)

server:
	go build -o $(SERVER) cmd/$(SERVER)/main.go

client:
	go build -o $(CLIENT) cmd/$(CLIENT)/main.go

clean:
	rm -f $(SERVER) $(CLIENT)
