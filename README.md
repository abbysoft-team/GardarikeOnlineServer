# Gardarike Online server
Gardarike Online is a new MMO strategy and rpg game. Currently it is in early stage of development. We have 2 developers working on the project, one at server side and one making frontend. What we want to achieve is to make gameplay unique but consistent between multiple platforms such as mobile devices and PC's. If you play on your mobile device, the game will be more like strategy economic building simulator. But once you switch platform to PC you will be able to enjoy plain old MMORPG genre expirience.

[![CircleCI](https://circleci.com/gh/abbysoft-team/GardarikeOnlineServer.svg?style=svg)](https://app.circleci.com/pipelines/github/abbysoft-team/GardarikeOnlineServer)

You can find client code here: https://github.com/abbysoft-team/GardarikeOnlineClient 

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

## LICENSE NOTICE
Feel free to use this code for non-profit goals. If you wan't to use it as part of commercial product contact us via contact@abbysoft.org. Usage without our (maintainers of this repo) permission is prohibited.
