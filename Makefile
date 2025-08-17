.PHONY: build install clean test help

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS = -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Build the binary
build:
	cd cmd/gwstool && go build -ldflags "$(LDFLAGS)" -o gwstool .

# Install to GOPATH/bin or GOBIN
install:
	cd cmd/gwstool && go install -ldflags "$(LDFLAGS)" .

# Clean build artifacts
clean:
	cd cmd/gwstool && rm -f gwstool

# Run tests
test:
	go test ./...

# Create example config directory and file
config-example:
	@mkdir -p ~/.config/gwstool
	@if [ ! -f ~/.config/gwstool/config ]; then \
		cp cmd/gwstool/config.example ~/.config/gwstool/config; \
		echo "Example config created at ~/.config/gwstool/config"; \
		echo "Please edit it with your credentials"; \
	else \
		echo "Config file already exists at ~/.config/gwstool/config"; \
	fi

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the gwstool binary"
	@echo "  install       - Install gwstool to GOPATH/bin"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  config-example - Create example config file"
	@echo "  help          - Show this help"
