package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type check struct {
	log *log.Logger
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return json.NewEncoder(w).Encode(struct{ Status string }{Status: "OK"})
}
