FROM golang as builder

WORKDIR /app

ADD . .

RUN go build

FROM debian

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

COPY --from=builder /app/bot .

ENTRYPOINT ["/app/bot"]
