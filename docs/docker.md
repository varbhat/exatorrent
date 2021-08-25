`exatorrent` can be run inside Podman / Docker Container 

## Podman
[Podman](https://github.com/containers/podman) is drop-in replacement for Docker . Podman has several advantages over Docker . `exatorrent` can be run in `podman` . If you prefer `podman` , you can alias it to `docker`

```bash
alias docker=podman
```

## Podman / Docker Image
Podman / Docker Images of `exatorrent` are officially available for `amd64` and `arm64` architectures . They are built on release of new version by Github Actions and are Hosted and available at [Github Container Registry](https://ghcr.io/varbhat/exatorrent) . Podman / Docker images of `exatorrent` can be pulled by 

```bash
docker pull ghcr.io/varbhat/exatorrent:latest
```

This pulls latest version of `exatorrent` Podman / Docker Image

## Building Podman / Docker Container Image
Podman / Docker Image of `exatorrent` can also be built in your machine if you intend to build . Following commands will build `exatorrent` Podman / Docker Image .

```bash
git clone https://github.com/varbhat/exatorrent
cd exatorrent
docker build -t "exatorrent" . 
```

## Usage
Podman / Docker Image of `exatorrent` can be run by following command 

```bash
docker run -p 5000:5000 -p 42069:42069 -v /path/to/directory:/exa/exadir ghcr.io/varbhat/exatorrent:latest
# docker run -p 5000:5000 -p 42069:42069 -v /path/to/directory:/exa/exadir exatorrent
```
5000 port is default port where Web Client and API are served . `42069` is default port for Torrent Client where Torrent Transfers occur. So, they need to be exposed . Also Refer [Usage](usage.md) .
