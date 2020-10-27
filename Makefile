build:
	go build -ldflags "-s -w" abbysoft/gardarike-online/cmd/gardarike

generate:
	mkdir -p ./rpc/generated
	protoc -Irpc/protocol --go_out=plugins=grpc:./rpc/generated rpc/protocol/server.proto
