BINARY_NAME=client
GO=go
GOFLAGS=
LDFLAGS=-ldflags="-s -w" 


all: build


build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .


run: build
	./$(BINARY_NAME)


clean:
	rm -f $(BINARY_NAME)


test:
	$(GO) test ./...


install:
	$(GO) mod tidy


info:
	$(GO) version
	$(GO) env


help:
	@echo "Makefile targets:"
	@echo "  all        - Build the project (default target)"
	@echo "  build      - Build the Go binary"
	@echo "  run        - Build and run the application"
	@echo "  clean      - Clean up the binary and temporary files"
	@echo "  test       - Run tests"
	@echo "  install    - Install dependencies"
	@echo "  info       - Show Go version and environment details"
