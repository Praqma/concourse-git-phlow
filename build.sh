#!/bin/bash

set -e -u -x

mkdir -p assets

export GOOS=linux
export GOARCH=amd64

go build -o assets/check check/main.go
go build -o assets/in in/main.go
go build -o assets/out out/main.go

chmod +x assets/check
chmod +x assets/in
chmod +x assets/out

docker build --no-cache -t groenborg/concourse-git-phlow:latest .
docker push groenborg/concourse-git-phlow:latest