## Build
FROM golang:1.19 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN cd cmd/cowboy && go build -o /cowboy

## Deploy
FROM amazonlinux:2022

WORKDIR /

USER nobody:nobody

COPY --from=build /cowboy /cowboy

EXPOSE 8000


ENTRYPOINT ["/cowboy"]