FROM golang:1.23-alpine AS builder
WORKDIR /src

COPY go.mod go.sum .
RUN go mod download
COPY . ./
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
	go build -o /usr/local/bin/temporal-server authorizer/server/main.go

FROM temporalio/auto-setup
WORKDIR /etc/temporal
COPY ./config/development.yaml development.yaml
COPY ./config/development.yaml config/development.yaml
COPY --from=builder  /usr/local/bin/temporal-server /usr/local/bin/temporal-server
CMD ["./entrypoint.sh","autosetup"]
