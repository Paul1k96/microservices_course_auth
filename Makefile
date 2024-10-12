LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate: generate/proto

generate/proto:
	@sh ./scripts/proto/generate.sh

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

go/lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

go/test:
	go test -v ./...

build:
	GOOS=linux GOARCH=amd64 go build -o bin/server cmd/main.go

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t paul1k888/microservice_course_auth:v0.0.1 .
	docker push paul1k888/microservice_course_auth:v0.0.1
