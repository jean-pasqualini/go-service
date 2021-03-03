// Package handlers contains the full set of handler functions and routes supported by the web api.
package handlers

import (
	authentication "github.com/jean-pasqualini/go-service/business/auth"
	"github.com/jean-pasqualini/go-service/business/mid"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger, auth  *authentication.Auth) http.Handler {
	app := web.NewApp(
		shutdown,
		mid.Logger(log),
		mid.Errors(log),
		mid.Panics(log),
		mid.Metrics(),
	)

	app.Handle(
		http.MethodGet,
		"/readiness",
		check{log: log, build: build}.readiness,
	)
	app.Handle(
		http.MethodGet,
		"/liveness",
		check{log: log, build: build}.liveness,
	)

	return app
}
