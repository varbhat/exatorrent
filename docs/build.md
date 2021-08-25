# Build Docs
This Documents how to build `exatorrent` from sources .

`exatorrent` is written in [Go](https://golang.org) . As such you need `go` to be installed in order to compile .

`exatorrent` is dependent on [sqlite3](https://www.sqlite.org) ( used as database ) and [libutp](https://github.com/anacrolix/go-libutp) ( used for uTP connections ) . As such , you also require C and C++ compilers to be installed in order to compile . [gcc](https://gcc.gnu.org/) is preferred but [clang](https://clang.llvm.org/) also compiles well . There is no need to install any `devel` packages or `sqlite` or headers in system to compile `exatorrent` , as Source comes bundeled with drivers [`crawshaw/sqlite`](https://github.com/crawshaw/sqlite) and [`libutp`](https://github.com/anacrolix/go-libutp) and C compiler is used to compile cgo code .


`exatorrent` comes with beautiful , small and performant Web Client . It is written in  [Svelte](https://svelte.dev/) + [TypeScript](https://www.typescriptlang.org/) . It gets bundled with amazing [esbuild](https://esbuild.github.io/) and the built web client is then embedded within the binary using Go's [embed](https://pkg.go.dev/embed) . [Node.js](https://nodejs.org/) is thus required to build Web Client . Note that Node.js is required only to build Web Client written in Svelte and not required thereafter (i.e exatorrent is not dependent on Node.js and Node.js in only required to build Web Client ).


For sake of Convenience although not necessary to build `exatorrent` , Build commands of `exatorrent` are written in [Makefile](../Makefile) . Install [`make`](https://www.gnu.org/software/make/) to execute make commands . Note that you can also manually type build commands instead of using `make` .


## Requirements

* [go](https://golang.org)
* [gcc](https://gcc.gnu.org/)
* [make](https://www.gnu.org/software/make/) ( to execute make commands )

Requirements to build Web Client :
* [Node.js](https://nodejs.org/) (`node` and `npm` must be available )

## Build
Since Web Client will be embedded within final binary , Web Client needs to be built first .

Web Client can be built by :

```bash
make web
```

After building web client , `exatorrent` can be built by :
```bash
make app
```

You can see built `exatorrent` in `build` directory .

If you don't have `make` installed , you can execute these commands manually to build exatorrent :
```bash
cd internal/web && npm install && npm run build
cd ../..
env CGO_ENABLED=1 go build -trimpath -buildmode=pie -ldflags '-extldflags "-static -s -w"' -o  build/exatorrent exatorrent.go
```
## Notes 
* See [Building Docker/Podman Images](./docker.md#building-podman--docker-container-image) if you want to build `exatorrent` Docker / Podman Images .
* If you don't want to build Web Client or want to skip building Web Client , you can do it by creating empty / dummy `index.html` file at `internal/web/build` directory ( Create `build` folder if it didn't exist ) .  Note that Web Client will not be available then .


