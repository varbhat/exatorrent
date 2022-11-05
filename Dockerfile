FROM docker.io/alpine:edge AS build
RUN apk add --no-cache git make musl-dev go nodejs npm gcc g++
WORKDIR /exa
ADD . /exa
RUN go mod tidy && make web && make app

FROM docker.io/alpine:edge
LABEL maintainer="varbhat"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="varbhat/exatorrent"
LABEL org.label-schema.description="self-hostable torrent client"
LABEL org.label-schema.url="https://github.com/varbhat/exatorrent"
LABEL org.label-schema.vcs-url="https://github.com/varbhat/exatorrent"
COPY --from=build --chown=1000:1000 /exa/build/exatorrent /exatorrent
USER 1000:1000
WORKDIR /exa
EXPOSE 5000
EXPOSE 42069
VOLUME /exa/exadir
ENTRYPOINT ["/exatorrent"]