include local.env

# Environment variables
LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.1.0
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.22.1
	GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@v0.5.0
	GOBIN=$(LOCAL_BIN) go install github.com/dmarkham/enumer@v1.5.9
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7
	@if [ ! -d bin/protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis bin/protogen/googleapis &&\
		mkdir -p  bin/protogen/google/ &&\
		mv bin/protogen/googleapis/google/api bin/protogen/google &&\
		rm -rf bin/protogen/googleapis ;\
	fi
	@if [ ! -d bin/protogen/validate ]; then \
		mkdir -p bin/protogen/validate &&\
		git clone https://github.com/envoyproxy/protoc-gen-validate bin/protogen/protoc-gen-validate &&\
		mv bin/protogen/protoc-gen-validate/validate/*.proto bin/protogen/validate &&\
		rm -rf bin/protogen/protoc-gen-validate ;\
	fi
	@if [ ! -d bin/protogen/protoc-gen-openapiv2 ]; then \
		mkdir -p bin/protogen/protoc-gen-openapiv2/options &&\
		git clone https://github.com/grpc-ecosystem/grpc-gateway bin/protogen/openapiv2 &&\
		mv bin/protogen/openapiv2/protoc-gen-openapiv2/options/*.proto bin/protogen/protoc-gen-openapiv2/options &&\
		rm -rf bin/protogen/openapiv2 ;\
	fi

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate: generate/proto generate/go

generate/go:
	@go generate ./...
	mkdir -p api/swagger
	$(LOCAL_BIN)/statik -src=api/swagger/ -include='*.css,*.html,*.js,*.json,*.png'

generate/proto:
	@sh ./scripts/proto/generate.sh

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

go/lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.yaml

go/test:
	go test -v ./...

go/test-cover:
	go test -coverprofile=coverage.out ./...

go/show-cover:
	go tool cover -html=coverage.out

local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

build:
	GOOS=linux GOARCH=amd64 go build -o bin/server cmd/main.go

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t paul1k888/microservice_course_auth:v0.0.1 .
	docker push paul1k888/microservice_course_auth:v0.0.1
