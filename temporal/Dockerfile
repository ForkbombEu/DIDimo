FROM golang:1.23-alpine AS builder
WORKDIR /src

COPY go.mod go.sum .
RUN go mod download
COPY . ./
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
	go build authorizer/server/main.go

CMD ["/src/run.sh"]
