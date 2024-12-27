.PHONY: resolver test-resolver test-infura

golang-deps:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62
	@go install github.com/fernandrone/linelint@0.0.6
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protobuf:
	protoc -I=./proto --go_out=./proto --go-grpc_out=./proto proto/*.proto

resolver:
	go run ./resolver

.PHONY: build_relayer build_resolver
build_relayer:
	@go build -o bin/relayer ./cmd/relayer/

build_resolver:
	@go build -o bin/resolver ./cmd/resolver/

.PHONY: build
build: build_relayer build_resolver

.PHONY: clean_build
clean_build: clean protobuf build

.PHONY: clean
clean: # for local usage
	@rm -rf bin/*
	@rm -rf proto/*.pb.go

.PHONY: check_lint
check_lint: # for local usage
	@golangci-lint run ./... 
	@linelint .

.PHONY: fix_lint
fix_lint: # for local usage
	@golangci-lint run --fix ./...
	@linelint -a .

.PHONY: test
test:
	@go test -v -race -count=1 ./...

.PHONY: testsum
testsum:
	@gotestsum --format testname -- -race -count=1 ./...

.PHONY: test_quick
test_quick:
	@go test -v ./...

make test-resolver:
	go test -v github.com/1inch/p2p-network/resolver ./resolver/...

make test-infura:
	go test -v github.com/1inch/p2p-network/resolver -testify.m=TestInfuraEndpoint

make test-integration:
	go test -v github.com/1inch/p2p-network/test 
