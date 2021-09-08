package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/varbhat/exatorrent/internal/core"
	"github.com/varbhat/exatorrent/internal/web"
)

func main() {
	core.Initialize()

	http.HandleFunc("/api/socket", core.SocketAPI)
	http.HandleFunc("/api/auth", core.AuthCheck)
	http.HandleFunc("/api/stream/", core.StreamFile)
	http.HandleFunc("/api/torrent/", core.TorrentServe)
	http.Handle("/", web.FrontEndHandler)

	if core.Flagconfig.UnixSocket != "" {
		// Run under specified Unix Socket Path
		core.Info.Println("Starting server at Path (Unix Socket)", core.Flagconfig.UnixSocket)
		usock, err := net.Listen("unix", core.Flagconfig.UnixSocket)
		if err != nil {
			core.Err.Fatalln("Failed listening", err)
		}
		go func() {
			if core.Flagconfig.TLSCertPath != "" && core.Flagconfig.TLSKeyPath != "" {
				core.Info.Println("Serving the HTTPS with TLS Cert ", core.Flagconfig.TLSCertPath, " and TLS Key", core.Flagconfig.TLSKeyPath)
				core.Err.Fatal(http.ServeTLS(usock, nil, core.Flagconfig.TLSCertPath, core.Flagconfig.TLSKeyPath))
			} else {
				core.Err.Fatal(http.Serve(usock, nil))
			}
		}()
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-c
		_ = usock.Close()
	} else {
		// Run under specified Port
		core.Info.Println("Starting server on", core.Flagconfig.ListenAddress)
		if core.Flagconfig.TLSCertPath != "" && core.Flagconfig.TLSKeyPath != "" {
			core.Info.Println("Serving the HTTPS with TLS Cert ", core.Flagconfig.TLSCertPath, " and TLS Key", core.Flagconfig.TLSKeyPath)
			core.Err.Fatal(http.ListenAndServeTLS(core.Flagconfig.ListenAddress, core.Flagconfig.TLSCertPath, core.Flagconfig.TLSKeyPath, nil))
		} else {
			core.Err.Fatal(http.ListenAndServe(core.Flagconfig.ListenAddress, nil))
		}
	}
}
