#!/bin/sh
set -e -u -x

cat cgp-version/version

mkdir -p $GOPATH/src/github.com/praqma
cp -R tollgate/ $GOPATH/src/github.com/praqma

# RESOLVE DEPENDENCIES - TEST AND PRODUCTION
cd $GOPATH/src/github.com/praqma/concourse-git-phlow

go get github.com/tools/godep
godep restore

#RUN WHOLE TEST SUITE IN VERBOSE MODE
#THE P FLAG ENSURES TESTS WILL RUN SEQUENTIALLY
#THEY WILL FAIL IN PARALLEL BECAUSE THE TESTFIXTURE CREATES CONFLICTING DIRECTORIES
#HOWEVER THIS IS NOT RELATED TO THE RESULTS OF THE TESTS
godep go test -v -p 1 ./...