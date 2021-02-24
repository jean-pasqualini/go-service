package handlers

import (
	"context"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"log"
	"math/rand"
	"net/http"
	"github.com/pkg/errors"
)

type check struct {
	log *log.Logger
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n % 2 == 0 {
		return errors.New("untrusted error")
	}

	return web.Respond(ctx, w, struct{ Status string }{Status: "OK"}, http.StatusOK)
}
