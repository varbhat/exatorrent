# Deploy Docs
This Documents some things on deployment of `exatorrent` .

# Run as Seperate User
Create new seperate user and group to run `exatorrent` . In this document , we assume that `exatorrent` is run as `exatorrent` user and `exatorrent` group but any name is fine . Run `exatorrent` as seperate user and group .

# Service
Using Services help to make sure that `exatorrent` will be running even after reboot and better logging by service Manager . Below are example Service files for `Runit` and `Systemd` .

## Runit
`run` file of `exatorrent` is as below .

```sh
#!/bin/sh
cd /path/to/directory
exec chpst -u exatorrent:exatorrent /path/to/exatorrent
```

## Systemd
Systemd unit file `exatorrent.service` is as follows . 

```
[Unit]
Description=exatorrent

[Service]
Type=simple
Restart=always
RestartSec=5s
User=exatorrent
Group=exatorrent
WorkingDirectory=/path/to/directory
ExecStart=/path/to/exatorrent

[Install]
WantedBy=multi-user.target
```

# Reverse Proxy
We recommend running `exatorrent` behind reverse proxy . Below are example configurations of Nginx , Haproxy and Caddy made to reverse proxy `exatorrent` . Please don't put Basic Auth or any kind of custom auth / request modification system as Authentication is handled by `exatorrent` itself . `/api/socket` is WebSocket endpoint and Reverse Proxying server must not hinder it .

## Nginx

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

## Haproxy

```HAProxy
frontend proxy
  #bind *:80
  bind *:443 ssl crt /path/to/tls/cert.pem
  default_backend exatorrent

backend exatorrent
  server exatorrent localhost:5000
```

## Caddy

```
https://the.domain.tld {
  reverse_proxy * localhost:5000
}
```
