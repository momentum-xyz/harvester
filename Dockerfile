# syntax=docker/dockerfile:1.3
FROM golang:1.18.2-alpine3.16 as build

WORKDIR /src

# Seperate step to allow docker layer caching
COPY go.* ./
RUN go mod download

COPY . ./
RUN go build ./cmd/harvester/...


FROM alpine:3.16.0 as runtime

WORKDIR /srv

COPY --from=build /src/harvester /srv/harvester
COPY --from=build /src/config.yaml /srv/config.yaml

CMD ["/srv/harvester"]