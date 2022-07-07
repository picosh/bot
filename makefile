build-irc:
	docker build -t neurosnap/erock-irc -f ./irc/Dockerfile ./irc
.PHONY: build

push-irc:
	docker push neurosnap/erock-irc
.PHONY: push

build-chat:
	docker build -t neurosnap/erock-chat -f ./chat/Dockerfile ./chat
.PHONY: build-chat

push-chat:
	docker push neurosnap/erock-chat
.PHONY: push-chat

build: build-irc build-chat
.PHONY: build

push: push-irc push-chat
.PHONY: push

bp: build push
.PHONY: bp
