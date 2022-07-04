build:
	docker build -t neurosnap/erock-irc .
.PHONY: build

push:
	docker push neurosnap/erock-irc
.PHONY: push

bp: build push
.PHONY: bp
