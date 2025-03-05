FROM golang:1.23-alpine AS builder
WORKDIR /src

COPY go.mod go.sum .
RUN go mod download
COPY . ./
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
	go build -o /usr/local/bin/temporal-server authorizer/server/main.go

FROM temporalio/auto-setup AS upstream

FROM alpine:latest AS temporal-server
WORKDIR /etc/temporal
RUN apk upgrade --no-cache \
    && apk add --no-cache ca-certificates tzdata bash curl
COPY --from=upstream /etc/temporal /etc/temporal
COPY --from=upstream /usr/local/bin/          /usr/local/bin/
COPY --from=builder  /usr/local/bin/temporal-server /usr/local/bin/temporal
CMD ["./entrypoint.sh","autosetup"]
