#!/bin/bash

set -ex

echo "Starting getting dependency"

go get -u $@ \
  github.com/golang/protobuf/protoc-gen-go \
  github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
  github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

brew install coreutils
brew install protobuf
brew install jq

echo "Successfully getting dependency"
