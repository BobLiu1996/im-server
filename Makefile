GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)



INTERNAL_PROTO_FILES=$(shell find ./internal/conf/protos -name *.proto)
API_PROTO_FILES=$(shell find ./api -name *.proto)


.PHONY: init
# init env
init:
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal/conf/protos \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal/conf \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api/protos/v1 \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api/v1 \
 	       --go-http_out=paths=source_relative:./api/v1 \
 	       --go-grpc_out=paths=source_relative:./api/v1 \
 	       --go-errors_out=paths=source_relative:./api/v1 \
	       --openapi_out=fq_schema_naming=true,default_response=false:./api \
	       --validate_out=paths=source_relative,lang=go:./api/v1 \
	       $(API_PROTO_FILES)

.PHONY: dao
dao:
	go run ./internal/data/gen/gen.go

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: wire
wire:
	cd cmd/im_server && wire && cd ../../

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
