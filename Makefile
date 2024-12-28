PROJECT_NAME 	?= didimo
ORGANIZATION 	?= forkbombeu
DESCRIPTION  	?= "SSI Compliance tool"
ROOT_DIR		?= $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BINARY_NAME 	?= $(PROJECT_NAME)
SUBDIRS			?= $(ROOT_DIR)/...
MAIN_SRC 		?= $(ROOT_DIR)/main.go
DOCS			?= $(ROOT_DIR)/docs
GOCMD 			?= go
GOBUILD			?= $(GOCMD) build
GOCLEAN			?= $(GOCMD) clean
GOTEST			?= $(GOCMD) test
GOTOOL			?= $(GOCMD) tool
GOGET			?= $(GOCMD) get
GOMOD			?= $(GOCMD) mod
GOINST			?= $(GOCMD) install
GOPATH 			?= $(shell $(GOCMD) env GOPATH)
GOBIN 			?= $(GOPATH)/bin

# Tools & Linters
GOLANGCI_LINT 	?= $(GOBIN)/golangci-lint
GOFUMPT 		?= $(GOBIN)/gofumpt
GOVULNCHECK 	?= $(GOBIN)/govulncheck
OVERMIND 		?= $(GOBIN)/overmind
AIR 			?= $(GOBIN)/air

# Submodules
AZC				= $(ROOT_DIR)/pocketbase/zencode/zenflows-crypto
WEBAPP			= $(ROOT_DIR)/webapp
WCZ				= $(WEBAPP)/client_zencode
BIN				= $(ROOT_DIR)/.bin
SLANGROOM 		= $(BIN)/slangroom-exec


DEPS = mise wget git tmux
K := $(foreach exec,$(DEPS),\
        $(if $(shell which $(exec)),some string,$(error "🥶 `$(exec)` not found in PATH please install it")))

all: help

$(BIN):
	@mkdir $(BIN)

$(SLANGROOM): | $(BIN)
	@wget https://github.com/dyne/slangroom-exec/releases/latest/download/slangroom-exec-$(shell uname)-$(shell uname -m) -O $(SLANGROOM)
	@chmod +x $(SLANGROOM)
	@@echo "slangroom-exec 😎 installed"

.git:
	@echo "🌱 Setup Git"
	@git init -q
	@git branch -m main
	@git add .

$(AZC): .git
	@rm -rf $@
	@git submodule --quiet add -f https://github.com/interfacerproject/zenflows-crypto $(AZC) && git submodule update --remote --init

$(WCZ): .git
	@rm -rf $@
	@cd $(WEBAPP) && git submodule --quiet add -f https://github.com/ForkbombEu/client_zencode $(WCZ) && git submodule update --remote --init

## Build
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

VERSION_STRATEGY 	= semver # git, semver, date
VERSION 			:= $(shell cat VERSION 2>/dev/null || echo "0.1.0")
GIT_COMMIT 			?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH 			?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_TIME 			?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_BY 			?= $(shell whoami)

version: ## ℹ️ Display version information
	@echo "$(CYAN)Version:	$(RESET)$(VERSION)"
	@echo "$(CYAN)Commit:		$(RESET)$(GIT_COMMIT)"
	@echo "$(CYAN)Branch:		$(RESET)$(GIT_BRANCH)"
	@echo "$(CYAN)Built:		$(RESET)$(BUILD_TIME)"
	@echo "$(CYAN)Built by: 	$(RESET)$(BUILD_BY)"
	@echo "$(CYAN)Go version:	$(RESET)$(shell $(GOCMD) version)"

dev: tools $(SLANGROOM) $(AZC) $(WCZ) ## 🚀 run in watch mode
	$(OVERMIND) start

test:
	$(GOTEST) $(SUBDIRS) -v

tidy:
	@$(GOMOD) tidy

build: tidy ## 📦 build the project into a binary
	@$(GOBUILD) -o $(BINARY_NAME) $(MAIN_SRC)
	@echo "📦 built"

doc: ## 📚 Serve documentation on localhost
	cd $(DOCS) && bun i
	cd $(DOCS) && bun run docs:dev --open

clean: ## 🧹 Clean files and caches
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@echo "🧹 cleaned"

tools:
	mise install
	@if [ ! -f "$(GOLANGCI_LINT)" ]; then \
		$(GOINST) github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@if [ ! -f "$(GOFUMPT)" ]; then \
		$(GOINST) mvdan.cc/gofumpt@latest; \
	fi
	@if [ ! -f "$(GOVULNCHECK)" ]; then \
		$(GOINST) golang.org/x/vuln/cmd/govulncheck@latest; \
	fi
	@if [ ! -f "$(OVERMIND)" ]; then \
		$(GOINST) github.com/DarthSim/overmind/v2@latest; \
	fi
	@if [ ! -f "$(AIR)" ]; then \
		$(GOINST) github.com/air-verse/air@latest; \
	fi

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
