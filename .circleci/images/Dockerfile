FROM ubuntu:16.04

WORKDIR /image

RUN apt-get update && apt-get -y upgrade && apt-get install -y libzmq5 libzmq3-dev wget git make unzip pkg-config
RUN wget https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz 
RUN tar -xvf go1.15.2.linux-amd64.tar.gz && mv go /usr/local
ENV GOROOT /usr/local/go
ENV PATH $GOROOT/bin:/root/go/bin:$PATH

RUN go get -u github.com/golang/protobuf/protoc-gen-go

RUN wget https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip
RUN unzip protoc-3.6.1-linux-x86_64.zip -d protoc3

RUN mv protoc3/bin/* /usr/local/bin/
RUN mv protoc3/include/* /usr/local/include/
