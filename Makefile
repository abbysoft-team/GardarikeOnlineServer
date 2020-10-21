build:
	go build -ldflags "-s -w"

generate:
	mkdir -p ./rpc/generated
	protoc -Irpc/protocol --go_out=plugins=grpc:./rpc/generated rpc/protocol/server.proto
