FROM golang:1.23-alpine AS builder
WORKDIR /src

COPY go.mod go.sum .
RUN go mod download
COPY . ./
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    go build -o didimo cmd/didimo/didimo.go

FROM debian:12-slim
RUN apt update
RUN apt install -y build-essential make bash curl git tmux wget ca-certificates
WORKDIR /app

COPY . ./
COPY --from=builder /src/didimo didimo

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ENV MISE_DATA_DIR="/mise"
ENV MISE_CONFIG_DIR="/mise"
ENV MISE_CACHE_DIR="/mise/cache"
ENV MISE_INSTALL_PATH="/usr/local/bin/mise"
ENV PATH="/mise/shims:$PATH"

RUN curl https://mise.run | sh
RUN mise trust
RUN mise i

ARG TARFILE=temporal_cli_latest_linux_amd64.tar.gz
RUN wget 'https://temporal.download/cli/archive/latest?platform=linux&arch=amd64' -O $TARFILE
RUN tar xf $TARFILE
RUN rm $TARFILE
RUN mv temporal /usr/local/bin
ENV POCKETBASE_URL=http://127.0.0.1:8090
COPY ./scripts/entry.sh /app/entry.sh
COPY ./.certs/mailpit+3.pem /usr/local/share/ca-certificates/mailpit.crt
RUN cat /usr/local/share/ca-certificates/mailpit.crt >> /etc/ssl/certs/ca-certificates.crt
WORKDIR webapp
RUN bun i
RUN bun run build

EXPOSE 8090
ENV PORT=5100

CMD ["/app/entry.sh"]
