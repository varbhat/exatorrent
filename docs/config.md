`exatorrent` is fully configurable. This page details configurations that are possible in `exatorrent`


## Options that are configurable during Runtime

There are options that can be configured during runtime and without restarting `exatorrent`. It can also be configured through Web Client (or API). See [EngConfig](https://github.com/varbhat/exatorrent/blob/main/internal/core/vars.go#L352) for more details about the engconfig.

When exatorrent gets started, it checks whether engconfig.json file is present in config directory (which is subdirectory of main directory of exatorrent) and if json is valid configuration, it gets applied. also, when you change the configuration during runtime, the json file gets updated.  You can generate sample engconfig(so that you can modify it to set value you want it to have) by passing flag `-engc` while starting the program.

## Options that are not configurable during runtime

Many of configuration of torrent engine `anacrolix/torrent` are only applied at start of engine and cannot be configurable during runtime. See [TorConfig](https://github.com/varbhat/exatorrent/blob/main/internal/core/vars.go#L42) for more details.

Note that most of torcconfig maps to [ClientConfig](https://github.com/anacrolix/torrent/blob/master/config.go#L23)

You can generate sample torcconfig(so that you can modify it to set value you want it to have) by passing flag `-torc` while starting the program. Note that if you don't want to configure , set it's value as `null`.

## Actions on Torrent Completion

`exatorrent` can listen to completion of torrent and call Hook on Completion. Hook is just a HTTP POST Request containing Infohash, Name, Completed Time of Completed Torrent sent to configured URL.
`listencompletion` of engconfig.json specifies whether the torrent must be listened for completion.
`hookposturl` of engconfig.json specifies URL where the Hook HTTP request must be posted.
`notifyoncomplete` of engconfig.json specifies whether the connected user( and owner of torrent) must be notified of completion via API/Web-Interface.
