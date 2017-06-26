#!/bin/bash

set -e -u -x
mkdir -p assets

export GOOS=linux
export GOARCH=amd64

VERSION=3

godep go build  -ldflags "-X github.com/praqma/concourse-git-phlow/repo.Version=`echo $VERSION`" -o assets/check check/check.go
godep go build -ldflags "-X github.com/praqma/concourse-git-phlow/repo.Version=`echo $VERSION`" -o assets/in in/in.go
godep go build -ldflags "-X github.com/praqma/concourse-git-phlow/repo.Version=`echo $VERSION`" -o assets/out out/out.go

chmod +x assets/check
chmod +x assets/in
chmod +x assets/out

docker build --no-cache -t groenborg/concourse-git-phlow:$1 .
docker push groenborg/concourse-git-phlow:$1


