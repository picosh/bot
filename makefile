DOCKER_CMD?=docker
REGISTRY?=registry.erock.io

setup:
	$(DOCKER_CMD) tag erock-bot $(REGISTRY)/erock-bot
.PHONY: setup

build-bot:
	$(DOCKER_CMD) build -t $(REGISTRY)/erock-bot -f ./bot/Dockerfile ./bot
.PHONY: build-bot

push-bot:
	$(DOCKER_CMD) push $(REGISTRY)/erock-bot
.PHONY: push-bot

build: build-bot
.PHONY: build

push: push-bot
.PHONY: push

bp: build push
.PHONY: bp
