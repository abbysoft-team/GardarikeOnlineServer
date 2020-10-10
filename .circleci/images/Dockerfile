FROM golang:1.15

RUN apt-get update && apt-get install -y protobuf-compiler
RUN go get -u github.com/golang/protobuf/protoc-gen-go
