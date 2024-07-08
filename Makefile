APP_NAME=trackdocs
OUTDIR=$(PWD)/build
GO=go

.PHONY: help build clean run stop logs wire proto doc

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

run: ## Runs all the containers
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 COMPOSE_PROJECT_NAME=trackdocs docker-compose up -d --build

stop: ## Stops all the containers
	@docker-compose down

wire: ## Generate wire files
	@wire ./cmd/api

setup-dev: ## Setup development environment
	@docker network create --subnet=172.28.0.0/16 --ip-range=172.28.5.0/24 --gateway=172.28.5.254 trackdocs || true
