.PHONY: build worker1 worker2 worker3 lint

BINARY_NAME=mygrep
GO=go

build:
	$(GO) build -o $(BINARY_NAME) .

worker1: build
	./$(BINARY_NAME) --server --addr :9001

worker2: build
	./$(BINARY_NAME) --server --addr :9002

worker3: build
	./$(BINARY_NAME) --server --addr :9003

lint:
	$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2 run
