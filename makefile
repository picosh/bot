DOCKER_CMD?=podman
REGISTRY?=registry.erock.io

setup:
	$(DOCKER_CMD) tag erock-irc $(REGISTRY)/erock-irc
	$(DOCKER_CMD) tag erock-chat $(REGISTRY)/erock-chat
	$(DOCKER_CMD) tag erock-bot $(REGISTRY)/erock-bot
.PHONY: setup

build-irc:
	$(DOCKER_CMD) build -t erock-irc -f ./irc/Dockerfile ./irc
.PHONY: build

push-irc:
	$(DOCKER_CMD) push $(REGISTRY)/erock-irc
.PHONY: push

build-chat:
	$(DOCKER_CMD) build -t erock-chat -f ./chat/Dockerfile ./chat
.PHONY: build-chat

push-chat:
	$(DOCKER_CMD) push $(REGISTRY)/erock-chat
.PHONY: push-chat

build-bot:
	$(DOCKER_CMD) build -t erock-bot -f ./bot/Dockerfile ./bot
.PHONY: build-bot

push-bot:
	$(DOCKER_CMD) push $(REGISTRY)/erock-bot
.PHONY: push-bot

build: build-irc build-chat build-bot
.PHONY: build

push: push-irc push-chat push-bot
.PHONY: push

bp: build push
.PHONY: bp
