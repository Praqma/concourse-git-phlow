#!/bin/bash
set -e -u -x

C_PATH=$(pwd)

VERSION=$(cat cgp-version/version)

mkdir -p $GOPATH/src/github.com/praqma
cp -R concourse-git-phlow/ $GOPATH/src/github.com/praqma

# RESOLVE DEPENDENCIES - TEST AND PRODUCTION
cd $GOPATH/src/github.com/praqma/concourse-git-phlow
go get github.com/tools/godep
godep restore

export GOOS=linux
export GOARCH=amd64


godep go build  -ldflags "-X github.com/praqma/concourse-git-phlow/repo.Version=`echo $VERSION`" -o $C_PATH/concourse-git-phlow/assets/check check/check.go
godep go build -ldflags "-X github.com/praqma/concourse-git-phlow/repo.Version=`echo $VERSION`" -o $C_PATH/concourse-git-phlow/assets/in in/in.go
godep go build -ldflags "-X github.com/praqma/concourse-git-phlow/repo.Version=`echo $VERSION`" -o $C_PATH/concourse-git-phlow/assets/out out/out.go

chmod +x $C_PATH/concourse-git-phlow/assets/check
chmod +x $C_PATH/concourse-git-phlow/assets/in
chmod +x $C_PATH/concourse-git-phlow/assets/out

cp -R $C_PATH/concourse-git-phlow/* $C_PATH/concourse-git-phlow-artifacts/