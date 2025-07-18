.PHONY: default simulate master worker client build proto

default: simulate

simulate:
	go run ./cmd/simulate.go

master:
	go run ./cmd/master

worker:
	go run ./cmd/worker

# client:
# 	go run ./cmd/client

build:
	@mkdir -p ./build/
	go build -o ./build/master ./cmd/master
	go build -o ./build/worker ./cmd/worker


proto:
	@mkdir -p ./internal/protocol/
	protoc --go_out=./internal/protocol/ --go-grpc_out=./internal/protocol/ ./internal/proto/*

clean:
	@rm -rf ./build/
	@rm -rf ./internal/protocol/