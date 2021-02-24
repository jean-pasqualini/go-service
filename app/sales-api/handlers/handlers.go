// Package handlers contains the full set of handler functions and routes supported by the web api.
package handlers

import (
	"github.com/jean-pasqualini/go-service/foundation/web"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger) http.Handler {
	app := web.NewApp(
		shutdown,
		func(handler web.Handler) web.Handler {
			return handler
		},
	)

	app.Handle(http.MethodGet, "/readiness", check{log: log}.readiness)

	return app
}
