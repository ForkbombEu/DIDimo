FROM golang:1.23-alpine AS builder
WORKDIR /src

COPY go.mod go.sum .
RUN go mod download
COPY . ./
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache

ARG DOCKERIZE_VERSION=v0.9.2
RUN  --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    go install github.com/jwilder/dockerize@${DOCKERIZE_VERSION}
RUN cp $(which dockerize) /usr/local/bin/dockerize

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
	go build -o /usr/local/bin/temporal authorizer/server/main.go

RUN mkdir -p      /etc/temporal \
    && mv config /etc/temporal/  \
    && mv *.sh  /etc/temporal/
RUN apk add git
RUN git clone --depth 1 https://github.com/temporalio/temporal/

FROM alpine:latest AS temporal-server
WORKDIR /etc/temporal
RUN apk upgrade --no-cache \
    && apk add --no-cache ca-certificates tzdata bash curl
COPY --from=builder /usr/local/bin/dockerize /usr/local/bin/
COPY --from=builder /usr/local/bin/temporal /usr/local/bin/
RUN        mkdir -p /etc/temporal /etc/temporal/schema
COPY --from=builder /etc/temporal /etc/temporal
COPY --from=builder /src/temporal/schema /etc/temporal/schema
COPY --from=builder /src/temporal/docker/config_template.yaml /etc/temporal/config/
CMD ["./entrypoint.sh"]
