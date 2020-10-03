build:
	CGO_ENABLED=0 go build -ldflags "-s -w"

generate:
	protoc -Irpc/protocol --go_out=plugins=grpc:./rpc/generated rpc/protocol/server.proto
