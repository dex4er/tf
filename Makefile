NAME := tf

ifneq (,$(wildcard .env))
	include .env
  export
endif

AWK = awk
DOCKER = docker
ECHO = echo
GO = go
GORELEASER = goreleaser
INSTALL = install
MAKE = make
PRINTF = printf
RM = rm
SORT = sort

ifeq ($(GOOS),windows)
EXE = .exe
else ifeq ($(shell $(GO) env GOOS),windows)
EXE = .exe
endif

BIN = $(NAME)$(EXE)

ifeq ($(OS),Windows_NT)
ifneq (,$(LOCALAPPDATA))
BINDIR = $(LOCALAPPDATA)\Microsoft\WindowsApps
else
BINDIR = C:\Windows\System32
endif
else
ifneq (,$(wildcard $(HOME)/.local/bin))
BINDIR = $(HOME)/.local/bin
else ifneq (,$(wildcard $(HOME)/bin))
BINDIR = $(HOME)/bin
else
BINDIR = /usr/local/bin
endif
endif

VERSION = $(shell ( git describe --tags --exact-match 2>/dev/null || ( git describe --tags 2>/dev/null || echo "0.0.0-0-g$$(git rev-parse --short=8 HEAD)" ) | sed 's/-[0-9][0-9]*-g/-SNAPSHOT-/') | sed 's/^v//' )
REVISION = $(shell git rev-parse HEAD)
BUILDDATE = $(shell TZ=GMT date '+%Y-%m-%dT%R:%SZ')

CGO_ENABLED = 0
export CGO_ENABLED

.PHONY: help
help:
	@echo Targets:
	@$(AWK) 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9._-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | $(SORT)

.PHONY: build
build: ## Build app binary for single target
	$(call print-target)
	$(GO) build -trimpath -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: goreleaser
goreleaser: ## Build app binary for all targets
	$(call print-target)
	$(GORELEASER) release --auto-snapshot --clean --skip-publish

$(BIN):
	@$(MAKE) build

.PHONY: install
install: ## Build and install app binary
install: $(BIN)
	$(call print-target)
	$(INSTALL) $(BIN) $(BINDIR)

.PHONY: uninstall
uninstall: ## Uninstall app binary
uninstall:
	$(RM) -f $(BINDIR)/$(BIN)

.PHONY: download
download: ## Download Go modules
	$(call print-target)
	$(GO) mod download

.PHONY: tidy
tidy: ## Tidy Go modules
	$(call print-target)
	$(GO) mod tidy

.PHONY: upgrade
upgrade: ## Upgrade Go modules
	$(call print-target)
	$(GO) get -u

.PHONY: clean
clean: ## Clean working directory
	$(call print-target)
	$(RM) -f $(BIN)
	$(RM) -rf dist

.PHONY: version
version: ## Show version
	@$(ECHO) "$(VERSION)"

.PHONY: revision
revision: ## Show revision
	@$(ECHO) "$(REVISION)"

.PHONY: builddate
builddate: ## Show build date
	@$(ECHO) "$(BUILDDATE)"

DOCKERFILE = Dockerfile
IMAGE_NAME = $(BIN)
LOCAL_REPO = localhost:5000/$(IMAGE_NAME)
DOCKER_REPO = localhost:5000/$(IMAGE_NAME)

ifeq ($(PROCESSOR_ARCHITECTURE),ARM64)
PLATFORM = linux/arm64
else ifeq ($(shell uname -m),arm64)
PLATFORM = linux/arm64
else ifeq ($(shell uname -m),aarch64)
PLATFORM = linux/arm64
else ifeq ($(findstring ARM64, $(shell uname -s)),ARM64)
PLATFORM = linux/arm64
else
PLATFORM = linux/amd64
endif

.PHONY: image
image: ## Build a local image without publishing artifacts.
	$(MAKE) build GOOS=linux
	$(call print-target)
	$(DOCKER) buildx build --file=$(DOCKERFILE) \
	--platform=$(PLATFORM) \
	--build-arg VERSION=$(VERSION) \
	--build-arg REVISION=$(REVISION) \
	--build-arg BUILDDATE=$(BUILDDATE) \
	--tag $(LOCAL_REPO) \
	--load \
	.

.PHONY: push
push: ## Publish to container registry.
	$(call print-target)
	$(DOCKER) tag $(LOCAL_REPO) $(DOCKER_REPO):v$(VERSION)-$(subst /,-,$(PLATFORM))
	$(DOCKER) push $(DOCKER_REPO):v$(VERSION)-$(subst /,-,$(PLATFORM))

.PHONY: test-image
test-image: ## Test local image
	$(call print-target)
	$(DOCKER) run --platform=$(PLATFORM) --rm -t $(LOCAL_REPO) -v

define print-target
	@$(PRINTF) "Executing target: \033[36m$@\033[0m\n"
endef
