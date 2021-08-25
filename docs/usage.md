## Usage
```bash
Usage of exatorrent:
 -addr    <addr> Listen Address (Default: ":5000")
 -admin   <user> Default admin username (Default Username: "adminuser" and Default Password: "adminpassword")
 -cert    <path> Path to TLS Certificate (Required for HTTPS)
 -dir     <path> exatorrent Directory (Default: "exadir")
 -engc    <opt>  Generate Custom Engine Configuration
 -key     <path> Path to TLS Key (Required for HTTPS)
 -psql    <opt>  Generate Sample Postgresql Connection URL
 -torc    <opt>  Generate Custom Torrent Client Configuration
 -unix    <path> Unix Socket Path
 -help    <opt>  Print this Help
 ```
 
 ### `-addr`
Listen Address of `exatorrent` . It specifies the TCP address for the `exatorrent` to listen on . It's of form `host:port` . The host must be a literal IP address, or a host name that can be resolved to IP addresses. The port must be a literal port number or a service name. If the host is a literal IPv6 address it must be enclosed in square brackets, as in `[2001:db8::1]:80` or `[fe80::1%zone]:80`. 

Default Listen Address is `:5000` . Open http://localhost:5000 or http://0.0.0.0:5000 or http://127.0.0.1:5000 to use `exatorrent` Web Client if exatorrent is listening on Default Address .

Valid Listen Addresses include `localhost:9999` , `0.0.0.0:3456` , `127.0.0.1:7777` , `[0:0:0:0:0:0:0:1]:5000` , `x.x.x.x:port` , `[x:x:x:x:x:x:x:x]:port` .

### `-admin`
Usernames of Users in `exatorrent`  can't be changed after User is created . It must be choosen while creating User and can't be changed later . Note that password can be changed anytime later .

On the first use of `exatorrent` , `exatorrent` creates Admin user with username `adminuser` and password `adminpassword` . Since this username `adminuser` can't be changed later on , you can customize it before first run itself by passing your desired username to `-admin` flag.

```bash
exatorrent -admin "mycustomadminusername"
```

You are advised to use custom username for default admin using this flag . You can always change password of default admin or other users later .

### `-cert`
If HTTPS needs to be served , Path to TLS Certificate must be specified in this flag . Note that `-cert` flag is also necessary along with this flag to serve HTTPS .

### `-key`
If HTTPS needs to be served , Path to TLS Key must be specified in this flag . Note that `-cert` flag is also necessary along with this flag to serve HTTPS .

### `-unix`
This flag specifies file path where Unix Socket must be served . This flag is alternative to `-addr` flag and works only on Operating Systems that support [Unix Sockets](https://en.wikipedia.org/wiki/Unix_domain_socket) .

### `-dir`
`exatorrent` always operates in specific directory and never leaves beyond that directory . It stores File downloaded through Torrents , Torrent Metadata Files , Sqlite3 Database files if any , Configuration files in this directory.

`-dir` flag specifies the directory where `exatorrent` must store the data .

### `-engc`
Writes Sample Runtime-Configurable Engine connection URL to `<exadirectory>/config/engconfig.json` . All Runtime-Configurable Settings are configured in this file . You can change it anytime in Client .

### `-torc`
Writes Sample Torrent Client connection URL to `<exadirectory>/config/clientconfig.json` . These settings cannot be configured during runtime and must only be configured before starting `exatorrent` .

### `-psql`
Writes Sample Postgresql connection URL to `<exadirectory>/config/psqlconfig.txt` . Instead of reading configuration URL from env variable `DATABASE_URL` , reads `DATABASE_URL` from file . Refer [Database Docs](database.md) for more details .


## Blocklist
Blocklist of PeerGuardian Text Lists (P2P) Format can be placed at `<exadirectory>/config/blocklist` to apply blocking based upon blocklist . You are also required to set `IPBlocklist` value of `clentconfig.json` to `true` . Note that Blocklist can't be changed during runtime .

## Rate Limiter
`UploadLimiterLimit` , `UploadLimiterBurst` , `DownloadLimiterLimit` , `DownloadLimiterBurst`  values of `clientconfig.json` can be used to Rate Limit `exatorrent` .
