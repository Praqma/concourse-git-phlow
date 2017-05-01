#!/bin/bash

set -e -u -x
mkdir -p assets

export GOOS=linux
export GOARCH=amd64

godep go build -o assets/check check/check.go
godep go build -o assets/in in/in.go
godep go build -o assets/out out/out.go

chmod +x assets/check
chmod +x assets/in
chmod +x assets/out

docker build --no-cache -t groenborg/concourse-git-phlow:$1 .
docker push groenborg/concourse-git-phlow:$1