package handlers

import (
	"context"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type check struct {
	build string
	log *log.Logger
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
	//	return web.NewRequestError(errors.New("trusted error"), http.StatusNotFound)
	}

	return web.Respond(ctx, w, struct{ Status string }{Status: "OK"}, http.StatusOK)
}

func (c check) liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	info := struct {
		Status string
		Build string
		Host string
		Pod string
		PodIP string
		Node string
		Namespace string
	} {
		Status: "up",
		Build: c.build,
		Host: host,
		Pod: os.Getenv("KUBERNETES_PODNAME"),
		PodIP: os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node: os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	return web.Respond(ctx, w, info, http.StatusOK)
}
