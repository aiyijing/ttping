VERSION := 1.0.0
BUILD_DIR := build

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: all clean

all: clean tidy $(PLATFORMS)

tidy:
	go mod tidy
clean:
	rm -rf $(BUILD_DIR)

$(PLATFORMS):
	$(eval GOOS := $(word 1,$(subst /, ,$@)))
	$(eval GOARCH := $(word 2,$(subst /, ,$@)))
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/ttping-$(GOOS)-$(GOARCH)

build: all
