# Makefile wrapper for go-task
# This project uses go-task (Taskfile.yml) as the primary build tool
# This Makefile provides familiar shortcuts for make users

.PHONY: help
help: ## Show this help message
	@echo "This project uses go-task. Run 'task' to see all available tasks."
	@echo ""
	@echo "Common shortcuts:"
	@echo "  make build       - Build the project"
	@echo "  make test        - Run tests"
	@echo "  make proto       - Generate proto files"
	@echo "  make run         - Run the server"
	@echo "  make clean       - Clean build artifacts"
	@echo ""
	@echo "For all tasks, run: task --list"

.PHONY: build
build: ## Build the project
	@task build

.PHONY: test
test: ## Run tests
	@task test

.PHONY: proto
proto: ## Generate proto files
	@task proto

.PHONY: run
run: ## Run the server
	@task run

.PHONY: clean
clean: ## Clean build artifacts
	@task clean

.PHONY: lint
lint: ## Run linters
	@task lint

.PHONY: fmt
fmt: ## Format code
	@task fmt

.PHONY: ci
ci: ## Run CI checks
	@task ci

.PHONY: install-tools
install-tools: ## Install development tools
	@task install-tools
