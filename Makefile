DOCKER_CMD?=docker
REGISTRY?=localhost:1338

fmt:
	go fmt ./...
.PHONY: fmt

setup:
	$(DOCKER_CMD) tag erock-bot $(REGISTRY)/erock-bot
.PHONY: setup

build:
	$(DOCKER_CMD) build -t $(REGISTRY)/erock-bot:latest .
.PHONY: build-bot

push:
	$(DOCKER_CMD) push $(REGISTRY)/erock-bot:latest
.PHONY: push-bot

bp: build push
.PHONY: bp
