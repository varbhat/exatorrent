//go:build noui
// +build noui

package web

import (
	"net/http"
)

// FrontEndHandler Provides Handler to Serve Frontend
var FrontEndHandler http.Handler

func init() {
	FrontEndHandler = http.NotFoundHandler()
}
