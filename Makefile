.PHONY: resolver

golang-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protobuf:
	protoc -I=./proto --go_out=./proto --go-grpc_out=./proto proto/*.proto

resolver:
	go run ./resolver
