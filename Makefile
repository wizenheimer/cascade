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
