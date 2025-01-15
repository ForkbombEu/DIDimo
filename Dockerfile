FROM golang:1.23-alpine AS builder
WORKDIR /src

COPY go.mod go.sum .
RUN go mod download
COPY . .
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    go build -o didimo cmd/didimo/didimo.go

FROM alpine:latest

RUN apk add --no-cache curl bash libstdc++

RUN curl -fsSL https://bun.sh/install | bash && \
    ln -s /root/.bun/bin/bun /usr/local/bin/bun
RUN curl -fsSL https://github.com/dyne/slangroom-exec/releases/latest/download/slangroom-exec-$(uname)-$(uname -m) -o /usr/bin/slangroom-exec && \
    chmod +x /usr/bin/slangroom-exec


WORKDIR /app

COPY --from=builder /src/didimo /usr/bin/didimo
COPY ./webapp/build/ /app/build
COPY ./scripts/entry.sh /app/entry.sh
COPY ./.certs/mailpit+3.pem /usr/local/share/ca-certificates/mailpit.crt
RUN cat /usr/local/share/ca-certificates/mailpit.crt >> /etc/ssl/certs/ca-certificates.crt
RUN apk --no-cache add --no-check-certificate ca-certificates && update-ca-certificates


EXPOSE 8090
ENV PORT=5100

CMD ["/app/entry.sh"]
