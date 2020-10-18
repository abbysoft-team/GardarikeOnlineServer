FROM ubuntu:16.04

WORKDIR /image

RUN apt-get update && apt-get -y upgrade && apt-get install -y protobuf-compiler libzmq5 libzmq3-dev wget git make
RUN wget https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz 
RUN tar -xvf go1.15.2.linux-amd64.tar.gz && mv go /usr/local
ENV GOROOT /usr/local/go
ENV PATH $GOROOT/bin:$PATH

RUN go get -u github.com/golang/protobuf/protoc-gen-go
