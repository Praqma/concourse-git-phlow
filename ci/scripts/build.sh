#!/bin/bash
set -e -u -x

go get github.com/tools/godep

cd concourse-git-phlow
godep restore

export GOOS=linux
export GOARCH=amd64

go build -o concourse-git-phlow/assets/check concourse-git-phlow/check/check.go
go build -o concourse-git-phlow/assets/in concourse-git-phlow/in/in.go
go build -o concourse-git-phlow/assets/out concourse-git-phlow/out/out.go

chmod +x concourse-git-phlow/assets/check
chmod +x concourse-git-phlow/assets/in
chmod +x concourse-git-phlow/assets/out

cp -R concourse-git-phlow/* concourse-git-phlow-artifacts/