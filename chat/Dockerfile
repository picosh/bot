FROM node:alpine AS builder

RUN apk update && apk add git

WORKDIR /app

RUN git clone https://git.sr.ht/~emersion/gamja

WORKDIR /app/gamja

RUN npm i --include=dev
RUN npm run build

FROM nginx:alpine

COPY ./nginx.conf /etc/nginx/nginx.conf

WORKDIR /srv

COPY --from=builder /app/gamja/dist .
COPY ./config.json /srv
