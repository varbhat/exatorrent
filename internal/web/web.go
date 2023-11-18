//go:build !noui
// +build !noui

package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed build/*
var webUI embed.FS

type webFS struct {
	Fs http.FileSystem
}

func (fs *webFS) Open(name string) (http.File, error) {
	f, err := fs.Fs.Open(name)
	if err != nil {
		return fs.Fs.Open("index.html")
	}
	return f, err
}

// FrontEndHandler Provides Handler to Serve Frontend
var FrontEndHandler http.Handler

func init() {
	contentStatic, _ := fs.Sub(fs.FS(webUI), "build")
	FrontEndHandler = http.FileServer(&webFS{Fs: http.FS(contentStatic)})
}
