package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type check struct {
	log *log.Logger
}

func (c check) readiness(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct { Status string }{ Status: "OK" })

	log.Println("OK")
}
