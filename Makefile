.PHONY: default simulate master worker client build

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
	go build -o ./build/master ./cmd/master 
	go build -o ./build/worker ./cmd/worker 