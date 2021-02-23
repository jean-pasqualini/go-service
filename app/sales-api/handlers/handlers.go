// Package handlers contains the full set of handler functions and routes supported by the web api.
package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
)

func API(build string, shutdown chan os.Signal, log *log.Logger) http.Handler {

	tm := httptreemux.NewContextMux()

	tm.Handle(http.MethodGet, "/readiness", check{log: log}.readiness)

	return tm
}
