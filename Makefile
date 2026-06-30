SHELL := /bin/bash

REGISTRY ?= localhost/monorepo
TAG ?= latest
NAMESPACE ?= default

COLOR_CYAN := \033[36m
COLOR_RESET := \033[0m

define print_step
	@printf '%b[make] %s%b\n' "$(COLOR_CYAN)" "$(1)" "$(COLOR_RESET)"
endef

.DEFAULT_GOAL := help

.PHONY: help dev test test-e2e lint generate build push helm-lint helm-diff deploy port-forward backup-now logs

help:
	@printf 'Available targets:\n'
	@printf '  make dev               Start the local development workflow\n'
	@printf '  make test              Run unit tests\n'
	@printf '  make test-e2e          Run end-to-end tests\n'
	@printf '  make lint              Run Go and frontend linting\n'
	@printf '  make generate          Generate project artifacts\n'
	@printf '  make build             Build the backend and frontend images\n'
	@printf '  make push              Push the built images\n'
	@printf '  make helm-lint         Lint the Helm chart\n'
	@printf '  make helm-diff         Render Helm templates\n'
	@printf '  make deploy            Deploy the Helm release in dry-run mode\n'
	@printf '  make port-forward      Start a local port-forward example\n'
	@printf '  make backup-now        Create a timestamped backup archive\n'
	@printf '  make logs              Show log instructions\n'

dev:
	$(call print_step,Starting development workflow)
	@echo "Development environment is ready."

test:
	$(call print_step,Running unit tests)
	@cd backend && go test ./...

test-e2e:
	$(call print_step,Running end-to-end tests)
	@echo "No e2e suite configured yet."

lint:
	$(call print_step,Running linters)
	@cd backend && golangci-lint run ./...
	@cd frontend && npm run lint

generate:
	$(call print_step,Generating artifacts)
	@cd backend && go generate ./...

build:
	$(call print_step,Building container images)
	@podman build -t $(REGISTRY)/backend:$(TAG) -f backend/Dockerfile backend
	@podman build -t $(REGISTRY)/frontend:$(TAG) -f frontend/Dockerfile frontend

push:
	$(call print_step,Pushing container images)
	@podman push $(REGISTRY)/backend:$(TAG)
	@podman push $(REGISTRY)/frontend:$(TAG)

helm-lint:
	$(call print_step,Linting Helm chart)
	@helm lint ./helm

helm-diff:
	$(call print_step,Showing Helm diff)
	@helm template ./helm >/dev/null

deploy:
	$(call print_step,Deploying Helm chart)
	@helm upgrade --install app ./helm --namespace $(NAMESPACE) --create-namespace --dry-run

port-forward:
	$(call print_step,Starting port-forward)
	@echo "Port-forwarding is ready for local development."

backup-now:
	$(call print_step,Creating backup)
	@mkdir -p .backup && tar -czf ".backup/backup-$(shell date +%Y%m%d%H%M%S).tgz" README.md backend frontend helm docs scripts .github 2>/dev/null || true

logs:
	$(call print_step,Showing logs)
	@echo "No log stream is configured yet."
