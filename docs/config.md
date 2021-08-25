`exatorrent` is fully configurable .  This page details configurations that are possible in `exatorrent`


## Options that are configurable during Runtime

There are options that can be configured during runtime and without restarting `exatorrent` . It can also be configured through Web Client (or API) . See [EngConfig](https://github.com/varbhat/exatorrent/blob/2420bebc9aa8f37cad959eeedc4b6abfdb4c8b28/internal/core/vars.go#L352)  for more details about the engconfig.

When exatorrent gets started, it checks whether engconfig.json file is present in config directory (which is subdirectory of main directory of exatorrent) and if json is valid configuration, it gets applied. also, when you change the configuration during runtime, the json file gets updated.  You can generate sample engconfig(so that you can modify it to set value you want it to have) by passing flag `-engc` while starting the program.

## Options that are not configurable during runtime

Many of configuration of torrent engine `anacrolix/torrent` are only applied at start of engine and cannot be configurable during runtime. See [TorConfig](https://github.com/varbhat/exatorrent/blob/aa8e587d64c6990dfaccdc4c9d415bf46d378593/internal/core/vars.go#L42) for more details.

Note that most of torcconfig maps to [ClientConfig](https://github.com/anacrolix/torrent/blob/v1.29.1/config.go#L23)

You can generate sample torcconfig(so that you can modify it to set value you want it to have) by passing flag `-torc` while starting the program. Note that if you don't want to configure , set it's value as `null` .
