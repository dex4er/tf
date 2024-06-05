NAME := tf

ifneq (,$(wildcard .env))
	include .env
  export
endif

AWK = awk
GO = go
GORELEASER = goreleaser
HEAD = head
INSTALL = install
PRINTF = printf
RM = rm
SORT = sort

ifeq ($(OS),Darwin)
SORT = gsort
else
SORT = sort
endif

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

.PHONY: help
help:
	@echo Targets:
	@$(AWK) 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9._-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | $(SORT)

.PHONY: build
build: ## Build app binary for single target
	$(call print-target)
	$(GORELEASER) build --clean --snapshot --single-target --output $(BIN)

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

.PHONY: test
test: ## Test app binary
test: $(BIN)
	$(call print-target)
	$(MAKE) -C tests test

define print-target
	@$(PRINTF) "Executing target: \033[36m$@\033[0m\n"
endef
