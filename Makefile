BINDIR       := $(CURDIR)/bin
INSTALL_PATH ?= /usr/local/bin
DIST_DIRS    := find * -type d -exec
PLATFORMS    := darwin/arm64 darwin/amd64 linux/arm64 linux/amd64
BINNAME      ?= mangathr

# go
PKG         := ./...
LDFLAGS     := -w -s

# Rebuild the binary if any of these files change
SRC := $(shell find . -type f -name '*.go' -print) go.mod go.sum

SHELL      = /usr/bin/env sh

# git
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
LDFLAGS   += -X github.com/browningluke/mangathr/v2/internal/version.sha=${GIT_SHA}

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X github.com/browningluke/mangathr/v2/internal/version.version=${BINARY_VERSION}
endif

# Default target
.PHONY: all
all: build


# install dependencies (for docker caching)

deps:
	go mod download
	go mod verify


# build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	go build -trimpath -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) ./cmd/mangathr

# install

.PHONY: install
install: build
	@install "$(BINDIR)/$(BINNAME)" "$(INSTALL_PATH)/$(BINNAME)"


# Release

# Build for all platforms
cross-build: LDFLAGS += -extldflags "-static"
cross-build: $(PLATFORMS)
	@echo "Build complete for all platforms."

# Define the build rule for each platform
$(PLATFORMS):
	@platform=$@; \
	GOARCH=$$(echo $$platform | cut -d'/' -f2) ; \
	GOOS=$$(echo $$platform | cut -d'/' -f1) ; \
	echo "Building for $$GOOS/$$GOARCH"; \
	mkdir -p dist/$$GOOS-$$GOARCH; \
	GOARCH=$$GOARCH GOOS=$$GOOS go build -o dist/$$GOOS-$$GOARCH/ -trimpath -ldflags '$(LDFLAGS)' $(PKG)

.PHONY: dist
dist:
	( \
		cd dist && \
		$(DIST_DIRS) cp ../LICENSE {} \; && \
		$(DIST_DIRS) cp ../README.md {} \; && \
		$(DIST_DIRS) tar -czf mangathr-${VERSION}-{}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r mangathr-${VERSION}-{}.zip {} \; \
	)


# Misc

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)' ./dist

.PHONY: info
info:
	@echo "Version:           ${VERSION}"
	@echo "Git Tag:           ${GIT_TAG}"
	@echo "Git Commit:        ${GIT_COMMIT}"
