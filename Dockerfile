# SPDX-License-Identifier: AGPL-3.0-or-later
# Copyright (C) 2022-2023 Dyne.org foundation <foundation@dyne.org>.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

ARG GOVER=1.20
ARG ALP_VER=3.18.0
FROM golang:$GOVER-alpine AS builder
ENV GONOPROXY=

RUN apk add --no-cache git

WORKDIR /app

COPY ./admin/go.mod ./admin/go.sum ./
RUN go mod download && go mod verify
COPY ./admin/ .
RUN git init
RUN rmdir zencode/zenflows-crypto
RUN git submodule add https://github.com/interfacerproject/zenflows-crypto zencode/zenflows-crypto

RUN ENABLE_CGO=0 go build -o pb .

FROM alpine:$ALP_VER

ARG USER=app

RUN wget https://github.com/dyne/zenroom/releases/latest/download/zenroom -O zenroom \
        && chmod a+x zenroom \
        && cp zenroom /usr/local/bin/

RUN addgroup --system "$USER" && adduser --system --ingroup "$USER" "$USER" && \
      install -d -m 0755 -o "$USER" -g "$USER" /log

WORKDIR /app

COPY --from=builder /app/pb .

# TODO: problem with pb_data directory
# USER $USER

EXPOSE 8090

ENTRYPOINT ["./pb"]
CMD ["serve", "--http=0.0.0.0:8090", "--publicDir=/app/pb_public"]
