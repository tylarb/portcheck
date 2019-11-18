all: test cli server

CLI=portcheck
SERVER=portcheck-server
BUILD_DIR?=$(CURDIR)/bin
CNTR_BUILD_DIR=/out/bin
GOVER=1.13
GOFLAGS=-mod=vendor


test:
	podman run -t --rm -v $(CURDIR):/portcheck --workdir /portcheck -e GOFLAGS=$(GOFLAGS) golang:$(GOVER) make test-local

$(BUILD_DIR)/$(CLI):
	go build -o $(BUILD_DIR)/$(CLI) ./cmd/$(CLI)

cli:
	podman run -t --rm -v $(CURDIR):/portcheck -v $(BUILD_DIR):/out --workdir /portcheck -e GOFLAGS=$(GOFLAGS) golang:$(GOVER) make -e BUILD_DIR=/out cli-local

$(BUILD_DIR)/$(SERVER):
	go build -o $(BUILD_DIR)/$(SERVER) ./cmd/$(SERVER)

server:
	podman run -t --rm -v $(CURDIR):/portcheck -v $(BUILD_DIR):/out --workdir /portcheck -e GOFLAGS=$(GOFLAGS) golang:$(GOVER) make -e BUILD_DIR=/out server-local


local: test-local cli-local server-local

test-local:
	go test ./cmd/portcheck
	go test ./cmd/portcheck-server

cli-local: $(BUILD_DIR)/$(CLI)

server-local: $(BUILD_DIR)/$(SERVER)



clean:
	rm -rf $(BUILD_DIR)/*
