# ProjectXServer
Project XXX

[![CircleCI](https://circleci.com/gh/abbysoft-team/ProjectXServer.svg?style=svg)](https://app.circleci.com/pipelines/github/abbysoft-team/ProjectXServer)

## Install

You can get fresh release from the Circle CI pipeline. Circle CI is building for ubuntu 16.04 right now.

## Building

Go with modules support is required.

## Linux (Elementary OS Hera) 

Install dependencies
```
sudo apt-get install golang-goprotobuf-dev
sudo apt-get install protobuf-compiler
```

Run in project folder
```
make generate
go build .
```
