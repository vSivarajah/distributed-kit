build:
	@go build -o bin/distributed-kit

run: build
	@./bin/distributed-kit

test:
	go test -v ./... 