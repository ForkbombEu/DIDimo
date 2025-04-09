# SPDX-FileCopyrightText: 2025 Forkbomb BV
#
# SPDX-License-Identifier: AGPL-3.0-or-later

PROJECT_NAME 	?= didimo
ORGANIZATION 	?= forkbombeu
DESCRIPTION  	?= "SSI Compliance tool"
ROOT_DIR		?= $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BINARY_NAME 	?= $(PROJECT_NAME)
SUBDIRS			?= ./...
MAIN_SRC 		?= $(ROOT_DIR)/cmd/didimo/didimo.go
DATA			?= $(ROOT_DIR)/pb_data
WEBAPP			?= $(ROOT_DIR)/webapp
GO_SRC 			:= $(wildcard **/*.go)
GODIRS			:= ./pkg/... ./cmd/...
UI_SRC			:= $(shell find $(WEBAPP)/src -type f \( -name '*.svelte' -o -name '*.js' -o -name '*.ts' -o -name '*.css' \) ! -name '*.generated.ts' ! -path 'webapp/src/modules/i18n/paraglide/*')
DOCS			?= $(ROOT_DIR)/docs
GOCMD 			?= go
GOBUILD			?= $(GOCMD) build
GOCLEAN			?= $(GOCMD) clean
GOTEST			?= $(GOCMD) test
GOTOOL			?= $(GOCMD) tool
GOGET			?= $(GOCMD) get
GOFMT			?= $(GOCMD) fmt
GOMOD			?= $(GOCMD) mod
GOINST			?= $(GOCMD) install
GOGEN			?= $(GOCMD) generate
GOPATH 			?= $(shell $(GOCMD) env GOPATH)
GOBIN 			?= $(GOPATH)/bin
GOMOD_FILES 	:= go.mod go.sum
COVOUT			:= coverage.out

# Tools & Linters
REVIVE			?= $(GOBIN)/revive
GOVULNCHECK 	?= $(GOBIN)/govulncheck
HIVEMIND 		?= $(GOBIN)/hivemind
GOW				?= $(GOBIN)/gow
GOCOVERTREEMAP	?= $(GOBIN)/go-cover-treemap

# Submodules
WEBENV			= $(WEBAPP)/.env
BIN				= $(ROOT_DIR)/.bin
DEPS 			= mise git temporal
DEV_DEPS		= pre-commit
K 				:= $(foreach exec,$(DEPS), $(if $(shell which $(exec)),some string,$(error "🥶 `$(exec)` not found in PATH please install it")))

all: help
.PHONY: submodules version dev test lint tidy purge build docker doc clean tools help w devtools

$(BIN):
	@mkdir $(BIN)

submodules:
	git submodule update --recursive --init

## Hacking
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

$(WEBENV):
	cp $(WEBAPP)/.env.example $(WEBAPP)/.env

dev: $(WEBENV) tools devtools submodules ## 🚀 run in watch mode
	$(HIVEMIND) Procfile.dev

test: ## 🧪 run tests with coverage
	$(GOTEST) $(GODIRS) -v -cover

ifeq (test.p, $(firstword $(MAKECMDGOALS)))
  test_name := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  $(eval $(test_name):;@true)
endif
test.p: tools ## 🍷 watch tests and run on change for a certain folder
	$(GOW) test -run "^$(test_name)$$" $(GODIRS)

coverage: devtools # ☂️ run test and open code coverage report
	-$(GOTEST) $(GODIRS) -coverprofile=$(COVOUT)
	$(GOTOOL) cover -html=$(COVOUT)
	$(GOCOVERTREEMAP) -coverprofile $(COVOUT) > coverage.svg && open coverage.svg

lint: devtools ## 📑 lint rules checks
	$(GOVULNCHECK) $(SUBDIRS)
	$(REVIVE) $(GODIRS)

fmt: devtools ## 🗿 format rules checks
	$(GOFMT) $(GODIRS)

tidy: $(GOMOD_FILES)
	@$(GOMOD) tidy

purge: ## ⛔ Purge the database
	@echo "⛔ Purge the database"
	@rm -rf $(DATA)

## Deployment

$(BINARY_NAME): $(GO_SRC) tools tidy submodules $(WEBENV)
	@$(GOBUILD) -o $(BINARY_NAME) $(MAIN_SRC)

$(WEBAPP)/build: $(UI_SRC)
	@./$(BINARY_NAME) serve & \
	PID=$$!; \
	./scripts/wait-for-it.sh localhost:8090 --timeout=60; \
	cd $(WEBAPP) && bun i && bun run build; \
	kill $$PID;

$(BINARY_NAME)-ui: $(UI_SRC)
	@./$(BINARY_NAME) serve & \
	PID=$$!; \
	./scripts/wait-for-it.sh localhost:8090 --timeout=60; \
	cd $(WEBAPP) && bun i && bun run bin; \
	kill $$PID;

docker: ## 🐳 run docker with all the infrastructure services
	docker compose build --build-arg PUBLIC_POCKETBASE_URL="http://localhost:8090"
	docker compose up

## Misc

doc: ## 📚 Serve documentation on localhost with --host
	cd $(DOCS) && bun i
	cd $(DOCS) && bun run docs:dev --open --host

clean: ## 🧹 Clean files and caches
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME)-ui
	@rm -fr $(WEBAPP)/build
	@rm -f $(DOCS)/.vitepress/config.ts.timestamp*
	@rm -f $(COVOUT) coverage.svg
	@echo "🧹 cleaned"

generate: $(ROOT_DIR)/pkg/gen.go
	$(GOGEN) $(ROOT_DIR)/pkg/gen.go

devtools: generate
	@if [ ! -f "$(REVIVE)" ]; then \
		$(GOINST) github.com/mgechev/revive@latest; \
	fi
	@if [ ! -f "$(GOVULNCHECK)" ]; then \
		$(GOINST) golang.org/x/vuln/cmd/govulncheck@latest; \
	fi
	@if [ ! -f "$(GOCOVERTREEMAP)" ]; then \
		$(GOINST) github.com/nikolaydubina/go-cover-treemap@latest; \
	fi
	@if [ ! -f "$(GOW)" ]; then \
		$(GOINST) github.com/mitranim/gow@latest; \
	fi
	pre-commit install
	pre-commit autoupdate
	-$(foreach exec,$(DEV_DEPS), $(if $(shell which $(exec)),some string,$(error "🥶 `$(exec)` not found in PATH please install it")))

tools: generate
	mise install
	@if [ ! -f "$(HIVEMIND)" ]; then \
		$(GOINST) github.com/DarthSim/hivemind@latest; \
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

kill-pocketbase: ## 🔪 Kill any running PocketBase instance
	@echo "Killing any existing PocketBase instance..."
	@-lsof -ti:8090 -sTCP:LISTEN | xargs kill -9 2>/dev/null || true

remove-overmind: ## 🧹 Remove overmind.sock
	@rm -f ./.overmind.sock
