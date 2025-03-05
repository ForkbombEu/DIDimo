FROM golang:1.23-alpine AS builder
WORKDIR /src

COPY go.mod go.sum .
RUN go mod download
COPY . ./
ARG DOCKERIZE_VERSION=v0.9.2
RUN go install github.com/jwilder/dockerize@${DOCKERIZE_VERSION}
RUN cp $(which dockerize) /usr/local/bin/dockerize
RUN apk upgrade --no-cache
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    bash \
    curl
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
	go build authorizer/server/main.go

CMD ["/src/entrypoint.sh"]
