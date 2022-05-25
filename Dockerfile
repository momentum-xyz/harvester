FROM golang:1.17.3-alpine as build

COPY . /usr/src/code
WORKDIR /usr/src/code
RUN go build ./cmd/harvester/...

FROM alpine:latest as production-build

RUN apk add --update --no-cache supervisor && rm -rf /var/cache/apk/*

RUN mkdir /opt/code
COPY --from=build /usr/src/code/harvester /opt/code/harvester
COPY --from=build /usr/src/code/config.yaml /opt/code/config.yaml

ADD supervisord.conf /etc/supervisord.conf

# This command runs your application, comment out this line to compile only
CMD ["/usr/bin/supervisord","-n", "-c", "/etc/supervisord.conf"]

LABEL Name=pengine Version=0.0.1
