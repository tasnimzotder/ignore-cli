BINARY_NAME=ignore
VERSION=$(shell git describe --tags --always --dirty)
PLATFORMS=darwin linux windows
ARCHITECTURES=amd64 arm64

LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# exit if the version is empty or not a tag
ifeq ($(VERSION),)
$(error "Version is not set")
endif

all: test build build-all clean

test:
	go test -v ./...

clean:
	go clean
	rm -rf build

build:
	go build $(LDFLAGS) -o build/$(BINARY_NAME) main.go

build-all:
	$(foreach GOOS, $(PLATFORMS),\
		$(foreach GOARCH, $(ARCHITECTURES),\
			$(shell export GOOS=$(GOOS) GOARCH=$(GOARCH) && \
				go build ${LDFLAGS} -o build/$(BINARY_NAME)-$(GOOS)-$(GOARCH)-$(VERSION) main.go) \
		))

	# rename windows binary
	mv build/$(BINARY_NAME)-windows-amd64-$(VERSION) build/$(BINARY_NAME)-windows-amd64-$(VERSION).exe
	mv build/$(BINARY_NAME)-windows-arm64-$(VERSION) build/$(BINARY_NAME)-windows-arm64-$(VERSION).exe

