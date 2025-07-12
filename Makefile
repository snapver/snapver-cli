BINARY_NAME=snapver-cli
VERSION ?= $(shell git describe --tags --always --abbrev=0 2>/dev/null || echo "v0.0.0")
LDFLAGS = -ldflags "-X github.com/snapver/snapver-cli/cmd.Version=$(VERSION)"

build:
	@echo "⚙️ Building Go binary with version info..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

start: build
	@echo "⚙️ Starting Snapver..."
	./$(BINARY_NAME) start

version: build
	@echo "⚙️ Showing version info..."
	./$(BINARY_NAME) version

tag:
	@echo "📌 Creating new tag..."
	@read -p "Enter version (e.g., v1.0.1): " version; \
	git tag $$version && echo "✅ Tag $$version created"

.PHONY: build start version tag