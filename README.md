# *exatorrent*
## 🧲 self-hostable torrent client

![GitHub Repo stars](https://img.shields.io/github/stars/varbhat/exatorrent)
![Latest Release](https://img.shields.io/github/release/varbhat/exatorrent)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/varbhat/exatorrent)
![GitHub License](https://img.shields.io/github/license/varbhat/exatorrent?logoColor=violet)

![Linux](https://img.shields.io/badge/Linux-%23.svg?logo=linux&color=FCC624&logoColor=black)
![macOS](https://img.shields.io/badge/macOS-%23.svg?logo=apple&color=000000&logoColor=white)
![Windows](https://img.shields.io/badge/Windows-%23.svg?logo=windows&color=0078D6&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%23.svg?logo=docker&color=1D63ED&logoColor=white)
![Podman](https://img.shields.io/badge/podman-%23.svg?logo=podman&color=734392&logoColor=white)

<p><a href="docs/screenshots.md">Screenshots</a> &bull; <a href="https://github.com/varbhat/exatorrent/releases/latest">Releases</a> &bull; <a href="#features">Features</a> &bull; <a href="#usage"> Installation </a> &bull; <a href="docs/usage.md"> Usage</a> &bull; <a href="docs/docker.md">Docker</a> &bull; <a href="docs/build.md"> Build </a> &bull; <a href="LICENSE">License</a></p>
<hr>

* exatorrent is Elegant [BitTorrent](https://www.bittorrent.org/) Client written in [Go](https://go.dev/).
* It is Simple, easy to use, yet feature rich.
* It can be run locally or be hosted in Remote Server with good resources.
* It is Single Completely Statically Linked Binary with Zero External Dependencies.
* It is lightweight and light on resources.
* It comes with Beautiful Responsive Web Client written in Svelte and Typescript.
* Thanks to documented [WebSocket](https://datatracker.ietf.org/doc/html/rfc6455) [API](docs/API.md) of exatorrent, custom client can be created.
* It supports Single User Mode and Multi User Mode.
* Torrented Files are stored in local disk can be downloaded and streamed via HTTP/Browser/Media Players.

<hr>
<p align="center">
<img src="https://raw.githubusercontent.com/varbhat/exatorrent/assets/main.png" alt="exatorrent web client" width=400 height=550 />
  <p align="center"><a href="docs/screenshots.md">More Screenshots →</a></p>
</p>
<hr>

## Usage
Exatorrent can be operated using either of the following methods:
* **Releases:** You can download precompiled binary for your Operating System from [Releases](https://github.com/varbhat/exatorrent/releases/latest). Mark it as executable and run it. Refer [Usage](docs/usage.md).
  ```bash
  wget https://github.com/varbhat/exatorrent/releases/latest/download/exatorrent-linux-amd64
  chmod u+x ./exatorrent-linux-amd64
  ./exatorrent-linux-amd64
  ```
 * **Docker:** Docker images of exatorrent are also provided which lets exatorrent to be run in a Docker container. See [Docker Docs](docs/docker.md).
   ```bash
   docker pull ghcr.io/varbhat/exatorrent:latest
   docker run -p 5000:5000 -p 42069:42069 -v /path/to/directory:/exa/exadir ghcr.io/varbhat/exatorrent:latest
   ```
 * **Manual Build:** exatorrent is open source and can be built from sources. See [Build Docs](docs/build.md).
   ```bash
   make web && make app
   ```

#### Notes:
* Note that **Username** and **Password** of Default User created on first run are `adminuser` and `adminpassword` respectively.
* You can change Password later but Username of Account can't be changed after creation. Refer [Usage](docs/usage.md#-admin).
* [Github Actions](https://github.com/features/actions) is used to build and publish [Releases](https://github.com/varbhat/exatorrent/releases/latest) and [Docker Images](https://ghcr.io/varbhat/exatorrent) of exatorrent.


## Features
* Single Executable File with No Dependencies
* Small in Size
* Cross Platform
* Download (or Build ) Single Executable Binary and run. That's it
* Open and Stream Torrents in your Browser
* Add Torrents by Magnet or by Infohash or Torrent File
* Individual File Control (Start, Stop or Delete )
* Stop, Remove or Delete Torrent
* Persistent between Sessions
* Stop Torrent once SeedRatio is reached (Optional)
* Perform Actions on Torrent [Completion](docs/config.md#actions-on-torrent-completion) (Optional)
* Powered by [anacrolix/torrent](https://github.com/anacrolix/torrent)
* Download/Upload [Rate limiter](docs/usage.md#rate-limiter) (Optional)
* Apply [Blocklist](docs/usage.md#blocklist) (Optional)
* [Configurable](docs/config.md) via Config File but works fine with Zero Configuration
* Share Files by Unlocking Torrent or Lock Torrent (protect by Auth)  to prevent External Access
* Retrieve or Stream Files via HTTP
* Multi-Users with Authentication
* Auto Add Trackers to Torrent from TrackerList URL
* Auto Fetch Torrent Metainfo from Online/Local Metainfo Cache
* Download Directory as Zip or as Tarball
* Stream directly on Browser or [VLC](https://www.videolan.org/vlc/) or [mpv](https://mpv.io/) or other Media Players
* [Documented API](docs/API.md)
* Uses Sqlite3 (embedded database with no setup and no configuration) by Default for [Database](docs/database.md) but PostgreSQL can be used instead too

## License
[GPL-v3](LICENSE)
