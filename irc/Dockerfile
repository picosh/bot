FROM golang:1.18-alpine as builder

WORKDIR /app

RUN apk add --no-cache git sqlite-libs sqlite-dev build-base

RUN git clone https://github.com/emersion/soju

WORKDIR /app/soju

ENV CGO_ENABLED=1
ENV LDFLAGS="-linkmode external -extldflags -static"

RUN go build -ldflags "$LDFLAGS" -tags sqlite_omit_load_extension ./cmd/soju
RUN go build -ldflags "$LDFLAGS" -tags sqlite_omit_load_extension ./cmd/sojuctl

FROM alpine

WORKDIR /app

ADD ./soju.config .
COPY --from=builder /app/soju/soju .
COPY --from=builder /app/soju/sojuctl .

ENTRYPOINT ["/app/soju", "-config", "/app/soju.config"]
