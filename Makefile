SHELL = /bin/sh

ENV ?= ./docker-compose.env

include $(ENV)

OPENCTI_TOKEN ?= $(OPENCTI_ADMIN_TOKEN)
OPENCTI_URL ?= $(OPENCTI_BASE_URL)
GOCTI_REPO := .

GENERATOR := ./tools/gocti_type_generator

PYTHON_CMD := python3
VENV := ./.venv
VENV_ACTIVATE := . $(VENV)/bin/activate
PYTHON := $(VENV)/bin/python
PIP := $(VENV)/bin/pip

COMPOSE_FILE := ./docker-compose.yml
COMPOSE_ENV_FILE := ./docker-compose.env

all: fmt lint

.PHONY: fmt
fmt:
	gofumpt -l -w .
	wsl --fix ./...
	ruff format $(GENERATOR)

.PHONY: lint
lint:
	golangci-lint run --fix
	ruff check $(GENERATOR) --fix

.PHONY: generate
generate: .go-generate fmt lint
	$(MAKE) fmt

.PHONY: test
test:
	export OPENCTI_URL=$(OPENCTI_URL) && \
	export OPENCTI_TOKEN=$(OPENCTI_TOKEN) && \
	go test -failfast -race ./... -timeout 120s

.PHONY: start-opencti
start-opencti:
	sudo docker compose --file $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) up -d

.PHONY: stop-opencti
stop-opencti:
	sudo docker compose --file $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) down

.PHONY: clean
clean:
	go clean -cache -testcache -modcache
	ruff clean
	rm -rf $(GENERATOR)/build
	rm -rf $(GENERATOR)/gocti_type_generator.egg-info
	rm -rf $(VENV)

.PHONY: .create-venv
.create-venv:
	if [ ! -d $(VENV) ]; then $(PYTHON_CMD) -m venv $(VENV); fi

.PHONY: .install-generator
.install-generator: .create-venv
	$(VENV_ACTIVATE) && \
	$(PIP) install $(GENERATOR)

.PHONY: .go-generate
.go-generate: .install-generator
	$(VENV_ACTIVATE) && \
	export GOCTI_REPO=$(GOCTI_REPO) && \
	export OPENCTI_URL=$(OPENCTI_URL) && \
	export OPENCTI_TOKEN=$(OPENCTI_TOKEN) && \
	go generate ./...
