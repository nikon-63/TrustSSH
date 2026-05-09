.PHONY: help test-cli test-lambda test-terraform test-all

.DEFAULT_GOAL := help

help: ## Show this help and available targets
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_.-]+:.*## / {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test-cli: ## (Tests) Run CLI tests
	@echo "Running CLI tests"
	./helpers/cli-test.sh

test-lambda: ## (Tests) Run Lambda tests
	@echo "Running Lambda tests"
	./helpers/lambda-test.sh

test-terraform: ## (Tests) Run Terraform checks
	@echo "Running Terraform checks"
	./helpers/terraform-test.sh

test-all: test-terraform test-lambda test-cli ## (Tests) Run all tests
