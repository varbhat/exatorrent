exatorrent is [BitTorrent](https://www.bittorrent.org) Client . It is written in `go` . It is made to be hosted in Servers to turn them into Torrent Server / Seedbox . It can also be run locally ( it's totally worthy because `exatorrent` provides several features not found in other Torrent Clients ) as well .

exatorrent is small in size and comes with no bloat . It is compiled statically into Single Executable Binary file and is dependent on nothing . No libraries , No frameworks are required to run exatorrent . Just download  exatorrent which is single Binary Executable and run it . That's it . You can also build `exatorrent` from sources as `exatorrent` is open source . Thus , deploying / running `exatorrent` is very easy . 

`exatorrent` built binary executables are provided with each new [releases](https://github.com/varbhat/exatorrent/releases/latest) . Docker / Podman Images are also provided . 

`exatorrent` is cross platform and can run on all major platforms . It can be run on [Linux](https://www.kernel.org/) , [Android](https://www.android.com/) ( through [Termux](https://termux.com/) ) , [Macos](https://developer.apple.com/macos/) and [Windows](https://www.microsoft.com/en-in/windows) . Linux users can Download Executable Binaries from [Releases](https://github.com/varbhat/exatorrent/releases/latest) or use [Docker/Podman](./docker.md) or [Build](./build.md) from Sources . Android Termux Users can run the Linux Builds from [Releases](https://github.com/varbhat/exatorrent/releases/latest) . Windows Users are advised to use [Docker](https://docs.docker.com/docker-for-windows/install/) (refer [docker docs](./docker.md) ) or use [WSL2](https://docs.microsoft.com/en-us/windows/wsl/) or [Build](./build.md) from Sources. Macos users are advised to use [Docker](https://docs.docker.com/docker-for-mac/install/) (refer [docker docs](./docker.md) ) or [Build](./build.md) from Sources .

You can Open and Stream Torrents in your Browser using exatorrent . You can Add Torrents by Magnet or by Infohash or Torrent File . Each file in Torrent can be Started , Stopped  or Deleted . Torrents can be Started , Stopped , Removed or Deleted . And all these things persist between sessions , i.e even if you restart `exatorrent` , these things do persist .


If you want to , then `exatorrent` can stop Torrents on reaching certain seedratio (which you need to set). You can also apply Blocklist to block peers if you want to. You can also Rate Limit `exatorrent` .

`exatorrent` works perfectly fine without any configuration , but  if you need to configure `exatorrent` , `exatorrent` can be fully configured with configuration files .


Files downloaded via Torrenting in `exatorrent` can be retrieved or streamed or shared ( with auth protection of course ) . Directories can be retrieved as Zip / Tarballs .


`exatorrent` is multi-users (with admin users as administrators ) but also works fine in single user context . All users are authenticated to access `exatorrent` .


`exatorrent` has few niceties like Adding Trackers to Torrent from TrackerList URL/s (which you can configure ofcourse ) which increases Peers thus making Torrenting faster . Also , it can fetch Torrent metainfo from Online Cache (ex. [iTorrents.org](https://itorrents.org/) or Local Cache thus making fetching of Torrent metainfo faster .

`exatorrent` uses `sqlite` as it's [database](./database.md) but `postgresql` can be used as well.

Video / Audio files can be streamed directly in Browser . They can also streamed directly on [VLC](https://www.videolan.org/vlc/) or [mpv](https://mpv.io/) or other Media Players .

`exatorrent` has [documented API](./API.md) that enables you to create other Clients for `exatorrent` and integrate other services with `exatorrent` .
