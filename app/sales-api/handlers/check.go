package handlers

import (
	"context"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"net/http"
)

type check struct {
	log *log.Logger
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	/**
		var u user.User
		if err := web.Decode(r, &u); err != nil {
			return err
		}
	*/

	if n := rand.Intn(100); n % 2 == 0 {
		return web.NewRequestError(errors.New("trusted error"), http.StatusNotFound)
	}

	return web.Respond(ctx, w, struct{ Status string }{Status: "OK"}, http.StatusOK)
}
