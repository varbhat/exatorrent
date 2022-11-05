FROM docker.io/alpine:edge AS build
RUN apk add --no-cache git make musl-dev go nodejs npm gcc g++
WORKDIR /exa
ADD . /exa
RUN go mod tidy && make web && make app

FROM docker.io/alpine:edge
COPY --from=build --chown=1000:1000 /exa/build/exatorrent /exatorrent
USER 1000:1000
WORKDIR /exa
EXPOSE 5000
EXPOSE 42069
VOLUME /exa/exadir
ENTRYPOINT ["/exatorrent"]