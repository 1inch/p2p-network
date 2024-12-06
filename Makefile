.PHONY: resolver

golang-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protobuf:
	protoc -I=./proto --go_out=./proto --go-grpc_out=./proto proto/*.proto

resolver:
	go run ./resolver

.PHONY: build_relayer
build_relayer:
	@go build -o bin/relayer ./cmd/relayer/

.PHONY: build
build: build_relayer

.PHONY: clean_build
clean_build: clean protobuf build

.PHONY: clean
clean: # for local usage
	@rm -rf bin/*
	@rm -rf proto/*.pb.go