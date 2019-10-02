
#!/bin/bash

# example commands to run in order to update the api proto definitions.
GIT_TAG="v1.2.0" # change as needed
go get -d -u github.com/golang/protobuf/protoc-gen-go
git -C "$(go env GOPATH)"/src/github.com/golang/protobuf checkout $GIT_TAG
go install github.com/golang/protobuf/protoc-gen-go

protoc -I api/ --go_out=plugins=grpc:api api/api.proto
