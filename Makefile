default: simulate

simulate:
	go run "./cmd/simulate.go"

master:
	go run "./cmd/master.go"

worker:
	go run "./cmd/worker.go"

client:
	go run "./cmd/client.go"

build:
	go build "./cmd/master.go" -o master
	go build "./cmd/worker.go" -o worker