FROM ghcr.io/varbhat/void-container:musl AS build
RUN xbps-install -Syu || xbps-install -yu xbps
RUN xbps-install -yu
RUN xbps-install -Sy git curl bash make go nodejs gcc
WORKDIR /exa
ADD . /exa
RUN go mod tidy && make web && make app

FROM gcr.io/distroless/base
COPY --from=build --chown=1000:1000 /exa/build/exatorrent /exatorrent
USER 1000:1000
WORKDIR /exa
ENTRYPOINT ["/exatorrent"]
