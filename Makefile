.DEFAULT_GOAL := help
.PHONY: help

# detect the operating system
OSFLAG 				:=
ifneq ($(OS),Windows_NT)
	UNAME_S := $(shell uname -s)
	UNAME_M := $(shell uname -m)
	ifeq ($(UNAME_S),Linux)
		OSFLAG += LINUX
	endif
	ifeq ($(UNAME_S),Darwin)
		ifeq ($(UNAME_M),arm64)
			OSFLAG += arm64
		else
			OSFLAG += OSX
		endif
	endif
endif


help: ## 🛟  Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-7s\033[0m %s\n", $$1, $$2}'

# - Setup - #

setup_git: ## 📦 Setup the Git (if not already done)
	@echo "📦 Setup Git"
	@(git status > /dev/null 2>&1 && echo "Git already set up") || (echo "git init" && git init)
	@echo " "

setup_submodules: setup_git ## 📦 Setup the submodules
	@echo "📦 Setup the submodules"
	rm -rf admin/zencode/zenflows-crypto
	rm -rf webapp/zenflows-crypto
	cd admin && git submodule add -f https://github.com/interfacerproject/zenflows-crypto zencode/zenflows-crypto
	cd webapp && git submodule add -f https://github.com/interfacerproject/zenflows-crypto zenflows-crypto
	@echo " "

setup_zenroom: ## 📦 Setup zenroom
	@echo "📦 Setup zenroom"
	@if ! command -v zenroom &> /dev/null; then \
		echo -n "Need zenroom executables, do you want to download them (works only on Osx not arm64 and Linux)? [y/N] " && read ans && [ $${ans:-N} = y ]; \
		if [ $(OSFLAG) == "OSX" ]; then \
				wget -O /usr/local/bin/zencode-exec.command https://github.com/dyne/zenroom/releases/latest/download/zencode-exec.command; \
				wget -O /usr/local/bin/zenroom.command https://github.com/dyne/zenroom/releases/latest/download/zenroom.command; \
				ln -s /usr/local/bin/zenroom.command /usr/local/bin/zenroom; \
				ln -s /usr/local/bin/zencode-exec.command /usr/local/bin/zencode-exec; \
				chmod +x /usr/local/bin/zencode-exec; \
				chmod +x /usr/local/bin/zenroom; \
		fi; \
		if [ $(OSFLAG) == "LINUX" ]; then \
			wget -O /usr/local/bin/zencode-exec https://github.com/dyne/zenroom/releases/latest/download/zencode-exec; \
			wget -O /usr/local/bin/zenroom https://github.com/dyne/zenroom/releases/latest/download/zenroom; \
			chmod +x /usr/local/bin/zencode-exec; \
			chmod +x /usr/local/bin/zenroom; \
		fi; \
	else \
		echo "Zenroom executables already exists"; \
	fi
	@echo " "

setup_backend: ## 📦 Setup the frontend
	@echo "📦 Setup the backend"
	@if [ ! -d ./admin/pb_data ]; then \
    	mkdir ./admin/pb_data; \
	fi
	@cd admin && ./setup

	@echo " "

setup_frontend: ## 📦 Setup the frontend
	@echo "📦 Setup the frontend"
	if [ ! -f ./webapp/.env ]; then \
		cp ./webapp/.env.example ./webapp/.env; \
	fi
	cd webapp && pnpm i
	@echo " "

setup: setup_submodules setup_backend setup_zenroom setup_frontend ## 📦 Setup the project

# - Running - #

# sleep 2 is needed to wait for the admin server to start
# before the webapp tries to connect to it and generate the schema
# check the webapp/package.json for the predev and prebuild scripts

be: ## ⚙️ Run the backend
	./admin/pb serve

fe: ## ⚙️ Run the frontend
	sleep 2 && cd webapp && pnpm serve

fe_dev: ## ⚙️ Watch the frontend
	sleep 2 && cd webapp && pnpm dev

dev: ## ⚙️ Run the project in development mode
	$(MAKE) be fe_dev -j2

up: setup ## ⚙️ Run the project
	$(MAKE) be fe -j2

doc: ## 📚 Serve documentation on localhost
	npx -p docsify-cli docsify serve ./docs

definitions: ## ⚙️ Generate type definitions and schema
	cd webapp && pnpm definitions

# - Cleaning - #

remove_git: ## 🧹 Remove git
	@echo "🧹 Removing git"
	rm -rf .git
	@echo " "

clean_submodules: ## 🧹 Clean submodules
	@echo "🧹 Clean submodules"
	rm -rf admin/zencode/zenflows-crypto
	rm -rf webapp/zenflows-crypto
	@echo " "

clean_build: ## 🧹 Clean project build
	@echo "🧹 Clean project build"
	rm -f admin/pb
	rm -fr webapp/node_modules
	rm -f webapp/src/lib/pocketbase/types.ts
	rm -f webapp/src/lib/pocketbase/schema/db_schema.json
	rm -f webapp/src/lib/rbac/roles.ts
	rm -f webapp/src/lib/features/list.ts
	@echo " "

clean: clean_submodules clean_build ## 🧹 Clean the project

purge: ## ⛔ Purge the database
	@echo "⛔ Purge the database"
	rm -fr admin/pb_data/*
	@echo " "
