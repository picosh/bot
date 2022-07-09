build-irc:
	docker build -t neurosnap/erock-irc -f ./irc/Dockerfile ./irc
.PHONY: build

push-irc:
	docker push neurosnap/erock-irc
.PHONY: push

bp-irc: build-irc push-irc
.PHONY: bp-irc

build-chat:
	docker build -t neurosnap/erock-chat -f ./chat/Dockerfile ./chat
.PHONY: build-chat

push-chat:
	docker push neurosnap/erock-chat
.PHONY: push-chat

bp-chat: build-chat push-chat
.PHONY: bp-chat

build-bot:
	docker build -t neurosnap/erock-bot -f ./bot/Dockerfile ./bot
.PHONY: build-bot

push-bot:
	docker push neurosnap/erock-bot
.PHONY: push-bot

bp-bot: build-bot push-bot
.PHONY: bp-bot

bp: bp-irc bp-chat bp-bot
.PHONY: bp
