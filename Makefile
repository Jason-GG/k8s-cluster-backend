# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=cicd
BINARY_LINUX=$(BINARY_NAME)_linux

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v

clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME) $(BINARY_LINUX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

run-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v
	scp $(BINARY_LINUX) user@your-amazon-linux-machine-ip:~/$(BINARY_LINUX)
	ssh user@your-amazon-linux-machine-ip ~/$(BINARY_LINUX)

.PHONY: all build build-linux clean run run-linux
