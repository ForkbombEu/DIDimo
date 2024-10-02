# SPDX-FileCopyrightText: 2024 The Forkbomb Company
#
# SPDX-License-Identifier: AGPL-3.0-or-later

.DEFAULT_GOAL := help
.PHONY: help

help: ## 🛟  Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[36m⦿ %-7s\033[0m %s\n\n", $$1, $$2}'

# - Folder structure - #

ROOT_DIR= $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

ADMIN	= $(ROOT_DIR)/admin
AZC			= $(ADMIN)/zencode/zenflows-crypto
PB			= $(ADMIN)/pb
DATA		= $(ADMIN)/pb_data

WEBAPP	= $(ROOT_DIR)/webapp
WZC			= $(WEBAPP)/zenflows-crypto
WCZ			= $(WEBAPP)/client_zencode

BIN		= $(ROOT_DIR)/.bin
ZENROOM		= $(BIN)/zenroom
ZENCODE		= $(BIN)/zencode-exec

export PATH := $(BIN):$(PATH)

# - Dependency check - #

DEPS = pnpm git wget go npx
K := $(foreach exec,$(DEPS),\
        $(if $(shell which $(exec)),some string,$(error "🥶 `$(exec)` not found in PATH please install it")))

# - Operating system detection - #

ifneq ($(OS),Windows_NT)
	UNAME_S := $(shell uname -s)
	UNAME_M := $(shell uname -m)
	ifeq ($(UNAME_S),Linux)
		OSFLAG := LINUX
		ZENROOM_URL = https://github.com/dyne/zenroom/releases/latest/download/zenroom
		ZENCODE_URL = https://github.com/dyne/zenroom/releases/latest/download/zencode-exec
	endif
	ifeq ($(UNAME_S),Darwin)
		ifeq ($(UNAME_M),arm64)
			OSFLAG := arm64
		else
			OSFLAG := OSX
			ZENROOM_URL = https://github.com/dyne/zenroom/releases/latest/download/zenroom.command
			ZENCODE_URL = https://github.com/dyne/zenroom/releases/latest/download/zencode-exec.command
		endif
	endif
endif

# - Setup: Zenroom - #

$(BIN):
	@mkdir $(BIN)

ifneq ($(OSFLAG),arm64)

$(ZENROOM): | $(BIN)
	@wget -q -O $@ $(ZENROOM_URL)
	@chmod +x $@
	@echo "zenroom 😎 installed"

$(ZENCODE): | $(BIN)
	@wget -q -O $@ $(ZENCODE_URL) 
	@chmod +x $@
	@echo "zencode-exec 🤭 installed"

else

$(ZENROOM):
	@echo "For usage on Apple ARM processors, please compile [zenroom] manually"

$(ZENCODE):
	@echo "For usage on Apple ARM processors, please compile [zencode-exec] manually"

endif

# - Setup: GIT - #

STARTERS_CHECK := $(shell pwd | grep -q "/starters/saas/didimo" && echo true || echo false)

ifeq ($(STARTERS_CHECK),false)

.git:
	@echo "🌱 Setup Git"
	@git init -q
	@git branch -m main
	@git add .

else

.git:
	@echo "Skipping git setup"

endif

# - Setup: Submodules - #

$(AZC): .git
	@rm -rf $@
	@cd $(ADMIN) && git submodule --quiet add -f https://github.com/interfacerproject/zenflows-crypto zencode/zenflows-crypto && git submodule update --remote --init

$(WZC): .git
	@rm -rf $@
	@cd $(WEBAPP) && git submodule --quiet add -f https://github.com/interfacerproject/zenflows-crypto zenflows-crypto && git submodule update --remote --init

$(WCZ): .git
	@rm -rf $@
	@cd $(WEBAPP) && git submodule --quiet add -f https://github.com/ForkbombEu/client_zencode client_zencode && git submodule update --remote --init

# - Setup: Project - #

$(DATA):
	mkdir -p $(DATA)

$(PB): $(DATA)
	@echo "📦 Setup the backend"
	@cd $(ADMIN) && go mod tidy && go build

$(WEBAPP)/.env:
	@cp $(WEBAPP)/.env.example $(WEBAPP)/.env

setup_frontend: $(WEBAPP)/.env
	@echo "📦 Setup the frontend"
	cd $(WEBAPP) && pnpm i

setup: $(AZC) $(WZC) $(WCZ) $(ZENROOM) $(ZENCODE) $(PB) setup_frontend ## 📦 Setup the project

# - Running - #

# sleep 2 is needed to wait for the admin server to start
# before the webapp tries to connect to it and generate the schema
# check the webapp/package.json for the predev and prebuild scripts

be: ## ⚙️ Run the backend
	$(PB) serve

be_remote:  ## ⚙️ Run the backend from remote
	$(PB) serve --http=0.0.0.0:8090

fe: ## ⚙️ Run the frontend
	sleep 2 && cd webapp && pnpm serve

fe_dev: ## ⚙️ Watch the frontend
	sleep 2 && cd webapp && pnpm dev

dev: ## ⚙️ Run the project in development mode
	$(MAKE) be fe_dev -j2

up: setup ## 💄 Run the project
	$(MAKE) be fe -j2

up_remote: setup  ## 💄 Run the project from remote
	$(MAKE) be_remote fe -j2

doc: ## 📚 Serve documentation on localhost
	npx -p docsify-cli docsify serve ./docs

definitions: ## ⚙️ Generate type definitions and schema
	cd webapp && pnpm definitions

# - Cleaning - #

clean: ## 🧹 Clean the project
	@rm -rf $(AZC) $(WZC) $(WCZ) $(BIN) $(PB)
	@rm -fr webapp/node_modules
	@rm -f webapp/src/lib/pocketbase/types.ts
	@rm -f webapp/src/lib/pocketbase/schema/db_schema.json
	@rm -f webapp/src/lib/rbac/roles.ts
	@rm -f webapp/src/lib/features/list.ts
	@echo "The project is ✨ cleaned"

purge: ## ⛔ Purge the database
	@echo "⛔ Purge the database"
	@rm -rf $(DATA)
