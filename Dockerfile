FROM golang:1.18-alpine3.16

RUN mkdir -p /app/db

WORKDIR /app

RUN apk --update upgrade
RUN apk add git make scdoc sqlite build-base
# See http://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN git clone https://git.sr.ht/~emersion/soju

ADD ./soju.config /app/soju.config

WORKDIR /app/soju

RUN CGO_ENABLED=1 make
RUN CGO_ENABLED=1 make install

WORKDIR /app

CMD ["soju", "-config", "/app/soju.config"]
