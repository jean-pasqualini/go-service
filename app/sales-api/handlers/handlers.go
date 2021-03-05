// Package handlers contains the full set of handler functions and routes supported by the web api.
package handlers

import (
	authentication "github.com/jean-pasqualini/go-service/business/auth"
	"github.com/jean-pasqualini/go-service/business/data/user"
	"github.com/jean-pasqualini/go-service/business/mid"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger, db *sqlx.DB, auth *authentication.Auth) http.Handler {
	app := web.NewApp(
		shutdown,
		mid.Logger(log),
		mid.Errors(log),
		mid.Panics(log),
		mid.Metrics(),
	)

	// Register checkController endpoints.
	checkCtrl := checkController{log: log, build: build, db: db}
	app.Handle(http.MethodGet, "/readiness", checkCtrl.readiness)
	app.Handle(http.MethodGet, "/liveness", checkCtrl.liveness)

	// Register user management and authentication endpoints
	userCtrl := userController{user: user.New(log, db), auth: auth}
	app.Handle(http.MethodGet, "/users/:page/:rows", userCtrl.query, mid.Authenticate(auth), mid.Authorize(log, authentication.RoleAdmin))
	app.Handle(http.MethodGet, "/users/:id", userCtrl.queryByID, mid.Authenticate(auth))
	app.Handle(http.MethodGet, "/users/token/:kid", userCtrl.token)
	app.Handle(http.MethodPost, "/users", userCtrl.create, mid.Authenticate(auth), mid.Authorize(log, authentication.RoleAdmin))
	app.Handle(http.MethodPut, "/users/:id", userCtrl.update, mid.Authenticate(auth), mid.Authorize(log, authentication.RoleAdmin))
	app.Handle(http.MethodDelete, "/users/:id", userCtrl.delete, mid.Authenticate(auth), mid.Authorize(log, authentication.RoleAdmin))

	return app
}
