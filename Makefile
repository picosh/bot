DOCKER_CMD?=docker
REGISTRY?=ghcr.io/picosh/bot

fmt:
	go fmt ./...
.PHONY: fmt

setup:
	$(DOCKER_CMD) tag bot $(REGISTRY)/bot
.PHONY: setup

build:
	$(DOCKER_CMD) build -t $(REGISTRY)/bot:latest .
.PHONY: build-bot

push:
	$(DOCKER_CMD) push $(REGISTRY)/bot:latest
.PHONY: push-bot

bp: build push
.PHONY: bp
