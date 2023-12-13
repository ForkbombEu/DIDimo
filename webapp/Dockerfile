FROM node:lts-alpine AS builder

RUN apk add --no-cache libc6-compat git
RUN wget "https://github.com/pnpm/pnpm/releases/latest/download/pnpm-linuxstatic-x64" -O /bin/pnpm && \
    chmod +x /bin/pnpm

WORKDIR /build

COPY ./webapp/ .
RUN git init
RUN git submodule add https://github.com/interfacerproject/zenflows-crypto zenflows-crypto


EXPOSE $NODE_PORT
RUN pnpm i
CMD pnpm serve
