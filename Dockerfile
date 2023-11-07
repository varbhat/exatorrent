# syntax=docker/dockerfile:1

# Build the web ui from source
FROM docker.io/node:18 AS build-node
WORKDIR /exa
ADD internal/web /exa/internal/web
ADD Makefile /exa/
RUN make web

# Build the application from source
FROM docker.io/golang:1.21-alpine3.18 AS build-go
WORKDIR /exa

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
COPY --from=build-node /exa/internal/web/build /exa/internal/web/build
RUN apk add --no-cache make gcc g++ && \
    make app

# Deploy the application binary into a lean image
FROM alpine:3.18

LABEL maintainer="varbhat"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="varbhat/exatorrent"
LABEL org.label-schema.description="self-hostable torrent client"
LABEL org.label-schema.url="https://github.com/varbhat/exatorrent"
LABEL org.label-schema.vcs-url="https://github.com/varbhat/exatorrent"

COPY --from=build-go --chown=1000:1000 /exa/build/exatorrent /exatorrent

RUN apk add -U --upgrade --no-cache \
    ca-certificates

USER 1000:1000

WORKDIR /exa

RUN mkdir -p exadir 

EXPOSE 5000 42069

VOLUME /exa/exadir

ENTRYPOINT ["/exatorrent"]

CMD ["-dir", "exadir"]
