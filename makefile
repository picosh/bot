DOCKER_CMD?=docker
REGISTRY?=local.imgs.sh:8443

setup:
	$(DOCKER_CMD) tag erock-bot $(REGISTRY)/erock-bot
.PHONY: setup

build-bot:
	$(DOCKER_CMD) build -t $(REGISTRY)/erock-bot:latest -f ./bot/Dockerfile ./bot
.PHONY: build-bot

push-bot:
	$(DOCKER_CMD) push $(REGISTRY)/erock-bot:latest
.PHONY: push-bot

build: build-bot
.PHONY: build

push: push-bot
.PHONY: push

bp: build push
.PHONY: bp
