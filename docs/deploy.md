# Deploy Docs
This Documents some things on deployment of `exatorrent` .

# Run as Seperate User
Create new seperate user and group to run `exatorrent` . In this document , we assume that `exatorrent` is run as `exatorrent` user and `exatorrent` group but any name is fine . Run `exatorrent` as seperate user and group .

# Service
Using Services help to make sure that `exatorrent` will be running even after reboot and better logging by service Manager . Below are example Service files for `Runit` , `Systemd` , `Sysvinit` and `Openrc` . Modify them to suit your system .

## Runit
`run` file of `exatorrent` is as below .

```sh
#!/bin/sh
cd /path/to/directory
exec chpst -u exatorrent:exatorrent /path/to/exatorrent
```

## Systemd
Systemd unit file `exatorrent.service` is as follows :

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

## Sysvinit
Sysvinit script is as follows:

```sh
#!/bin/sh
### BEGIN INIT INFO
# Provides:          exatorrent
# Required-Start:    $remote_fs $network
# Required-Stop:     $remote_fs $network
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: starts exatorrent
# Description:       exatorrent is torrent client
### END INIT INFO

PATH=/bin:/usr/bin:/sbin:/usr/sbin
NAME=exatorrent
USER=exatorrent
DESC="exatorrent"
DAEMON_PATH=/path/to/exatorrent/bin/dir
PIDFILE=/run/$NAME.pid
LOGFILE=/var/log/$NAME.log
TIMEOUT=30
SCRIPTNAME=/etc/init.d/$NAME

case "$1" in
  start)
    echo "Starting $DESC"
    start-stop-daemon --start --background --oknodo --startas $DAEMON --chdir $DAEMON_PATH --chuid $USER --exec $DAEMON_PATH/$NAME -- 
    ;;
  stop)
    echo "Stoping $DESC"
    start-stop-daemon --stop --quiet --oknodo --retry=0/10/TERM/5/KILL/5 --exec $DAEMON_PATH/$NAME
    ;;
  status)
    start-stop-daemon --status --exec $DAEMON_PATH/$NAME && exit_status=$? || exit_status=$?
    case "$exit_status" in
        0)     echo "The '$DESC' is running."            ;;
        *)     echo "The '$DESC' is not running."            ;;
    esac
    ;;
  restart)
    echo "Restarting $DESC: "
    sh $0 stop
    sleep 20
    mv --backup=numbered $LOGFILE $LOGFILE.1
    sh $0 start
    ;;
esac
exit 0
```

## openrc

Openrc `run` file is as follows :

```sh
#!/sbin/openrc-run
description="exatorrent"

depend() {
	need net
	after firewall
}

start() {
	ebegin "Starting exatorrent"
	start-stop-daemon --start --background --chdir /path/to/exatorrent/dir --exec /path/to/exatorrent/dir/exatorrent --
	eend $?
}

stop() {
	local rv=0
	ebegin "Stopping exatorrent"
	start-stop-daemon --stop --quiet --oknodo --retry=0/10/TERM/5/KILL/5 --exec /path/to/exatorrent/dir/exatorrent 
	eend $?
}
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
