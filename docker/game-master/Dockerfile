## Build
FROM golang:1.19 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN cd cmd/game-master && go build -o /game-master

## Deploy
FROM amazonlinux:2022

WORKDIR /

USER nobody:nobody

COPY --from=build /game-master /game-master

EXPOSE 8000


ENTRYPOINT ["/game-master"]