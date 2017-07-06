#!/bin/sh
set -e -u -x

cat cgp-version/version

#Change some of this
mkdir -p $GOPATH/src/github.com/praqma/concourse-git-phlow
cp -R tollgate/* $GOPATH/src/github.com/praqma/concourse-git-phlow

# RESOLVE DEPENDENCIES - TEST AND PRODUCTION
cd $GOPATH/src/github.com/praqma/concourse-git-phlow

go get github.com/tools/godep
go get gopkg.in/zorkian/go-datadog-api.v2
godep restore

#RUN WHOLE TEST SUITE IN VERBOSE MODE
#THE P FLAG ENSURES TESTS WILL RUN SEQUENTIALLY
#THEY WILL FAIL IN PARALLEL BECAUSE THE TESTFIXTURE CREATES CONFLICTING DIRECTORIES
#HOWEVER THIS IS NOT RELATED TO THE RESULTS OF THE TESTS
godep go test -v -p 1 ./...