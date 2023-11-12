`exatorrent` can be run inside Docker Container

## Docker Image
Docker Images of `exatorrent` are officially available for `amd64`,`arm64` and `arm32` architectures. They are built on release of new version by Github Actions and are Hosted and available at [Github Container Registry](https://ghcr.io/varbhat/exatorrent). Alpine is used as Base Image of `exatorrent` image. Docker images of `exatorrent` can be pulled by

```bash
docker pull ghcr.io/varbhat/exatorrent:latest
```

This pulls latest version of `exatorrent` Docker Image. You can use [diun](https://github.com/crazy-max/diun) to keep exatorrent updated to latest release as when it rolls out.

## Docker Container Image
Docker Image of `exatorrent` can also be built in your machine if you intend to build. Following commands will build `exatorrent` Docker Image.

```bash
git clone https://github.com/varbhat/exatorrent
cd exatorrent
docker build -t "exatorrent".
```

## Usage
Docker Image of `exatorrent` can be run by following command

```bash
docker run -p 5000:5000 -p 42069:42069 -v /path/to/directory:/exa/exadir ghcr.io/varbhat/exatorrent:latest
# docker run -p 5000:5000 -p 42069:42069 -v /path/to/directory:/exa/exadir exatorrent
```
5000 port is default port where Web Client and API are served. `42069` is default port for Torrent Client where Torrent Transfers occur. So, they need to be exposed. Also Refer [Usage](usage.md).

You can use `--user` flag of docker to run `exatorrent` as other user.


# Deploy Docs

## Reverse Proxy
We recommend running `exatorrent` behind reverse proxy . Below are example configurations of Nginx , Haproxy and Caddy made to reverse proxy `exatorrent` . Please don't put Basic Auth or any kind of custom auth / request modification system as Authentication is handled by `exatorrent` itself . `/api/socket` is WebSocket endpoint and Reverse Proxying server must not hinder it .

### Nginx

```Nginx
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    ssl_certificate /path/to/tls/tls.crt;
    ssl_certificate_key /path/to/tls/tls.key;

    server_name the.domain.tld;

    location / {
        proxy_pass http://localhost:5000;
        # proxy_pass http://unix:/path/to/exatorrent/unixsocket;
    }

    location /api/socket {
        proxy_pass http://localhost:5000/api/socket;
        # proxy_pass http://unix:/path/to/exatorrent/unixsocket:/api/socket;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
    }

}
```

### Haproxy

```HAProxy
frontend proxy
  #bind *:80
  bind *:443 ssl crt /path/to/tls/cert.pem
  default_backend exatorrent

backend exatorrent
  server exatorrent localhost:5000
```

### Caddy

```
https://the.domain.tld {
  reverse_proxy * localhost:5000
}
```
