#!/usr/bin/env bash

# Get current directory
DIR="$(pwd)"

# Find all directories containing at least one protofile.
for dir in $(find "${DIR}/api/proto" -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq); do
    files=$(find "${dir}" -name '*.proto')

    protoc -I="${DIR}/api/proto" \
        --go_out=paths=source_relative:"${DIR}/pkg/proto/gen" \
        --plugin=protoc-gen-go=bin/protoc-gen-go \
        --go-grpc_out=paths=source_relative:"${DIR}/pkg/proto/gen" \
        --plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
        ${files}
done
