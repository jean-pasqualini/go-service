package main

import (
	"github.com/dimfeld/httptreemux/v5"
	//"errors"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {

	if 1 == 2 {
		errors.New("random error")
	}

	m := httptreemux.NewContextMux()
	m.Handle(
		http.MethodGet,
		"/test",
		nil,
	)

	return nil
}