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

WORKDIR /app

COPY --from=builder /src/didimo /usr/bin/didimo
COPY ./webapp/build/ /app/build
COPY ./scripts/entry.sh /app/entry.sh

EXPOSE 8090
ENV PORT=5100

CMD ["/app/entry.sh"]
