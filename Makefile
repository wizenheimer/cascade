.PHONY: init
init:
	@echo "== Initalizing Development Environment =="
	brew install go
	brew install node
	brew install pre-commit
	brew install golangci-lint
	brew upgrade golangci-lint

	@echo "== Installing Pre-Commit Hooks =="
	pre-commit install
	pre-commit autoupdate
	pre-commit install --install-hooks
	pre-commit install --hook-type commit-msg

.PHONY: dev
start:
	@echo "== Starting Development Environment =="
	docker compose up -d --build

.PHONY: stop
stop:
	@echo "== Stopping Development Environment =="
	docker compose down -v --remove-orphans

.PHONY: build
build:
	@echo "== Building Docker Image =="
	docker compose build

.PHONY: cli
cli:
	@echo "== Building CLI from Source =="
	cd src && go build -o ../cascade ./interface/cli && cd ..
