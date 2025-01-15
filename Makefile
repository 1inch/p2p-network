SHELL := /bin/bash

.PHONY: resolver test-resolver test-infura

golang-deps:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62
	@go install github.com/fernandrone/linelint@0.0.6
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/ethereum/go-ethereum/cmd/abigen@latest
	@go install go.uber.org/mock/mockgen@latest

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

.PHONY: generate_bindings
generate_bindings:
	@go generate -x ./contracts

.PHONY: generate_mocks
generate_mocks:
	@go generate -x -run="mockgen" ./...

.PHONY: clean_build
clean_build: clean protobuf build

.PHONY: clean
clean: # for local usage
	@rm -rf bin/*
	@rm -rf proto/*.pb.go
	@rm -rf contracts/bin/*

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

.PHONY: deploy_contract
deploy_contract:
	@echo "Deploying contract..."
	@go test -v -tags=deploy ./contracts -run ^TestDeployContract$

.PHONY: register_resolver
register_resolver:
	@go test -v -tags=deploy ./contracts -run ^TestRegisterResolver$

.PHONY: test_quick
test_quick:
	@go test -v ./...

test-resolver:
	go test -v github.com/1inch/p2p-network/resolver ./resolver/...

test-encryption:
	go test -v github.com/1inch/p2p-network/internal/encryption

test-infura:
	go test -v github.com/1inch/p2p-network/resolver -testify.m=TestInfuraEndpoint

.PHONY: start-anvil
start-anvil:
	@if ! command -v anvil &> /dev/null; then \
		echo "Anvil binary not found"; \
	elif ! pgrep -x "anvil" > /dev/null; then \
		echo "Starting anvil on port 8545"; \
		anvil & \
		timeout 5 sh -c 'until nc -z localhost 8545; do sleep 1; done' || (echo "Anvil failed to start." && exit 1); \
	else \
		echo "Anvil is already running"; \
	fi

.PHONY: stop-anvil
stop-anvil:
	@echo "Stopping Anvil..."
	@pids=`ps aux | grep 'anvil' | awk '{print $$2}'`; \
	if [ -n "$$pids" ]; then \
		echo "Found Anvil PIDs: $$pids"; \
		echo "$$pids" | xargs kill -9 && echo "Anvil processes killed."; \
	else \
		echo "No Anvil processes found."; \
	fi

test-integration:
	go test -v github.com/1inch/p2p-network/test 
