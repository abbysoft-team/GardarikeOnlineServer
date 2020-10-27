build:
	go build -ldflags "-s -w" -o gardarike-online abbysoft/gardarike-online/cmd/gardarike

generate:
	mkdir -p ./rpc/generated
	protoc -Irpc/protocol --go_out=plugins=grpc:./rpc/generated rpc/protocol/server.proto
