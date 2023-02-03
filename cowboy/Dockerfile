## Build
FROM golang:1.19 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN cd cmd/cowboy && go build -o /cowboy

## Deploy
FROM alpine:3.17

WORKDIR /

RUN adduser -D nonroot

USER nonroot:nonroot

COPY --from=build /cowboy /cowboy

EXPOSE 8000


ENTRYPOINT ["/cowboy"]