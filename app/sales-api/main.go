package main

import (
	"context"
	"expvar"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
Need to figure out timeouts for http service.
You might want to reset your DB_HOST env var during test tear down.
Service should start even without a DB running yet.
symbols in profiles:
- https://github.com/golang/go/issues/23376
- https://github.com/google/pprof/pull/366
*/

const NAMESPACE_CONF = "SALES"
const BUILD = "develop"

func main() {
	log := log.New(os.Stdout, "SALES : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}


func run(log *log.Logger) error {

	// =================================================================================================================
	// Configuration

	var cfg struct {
		conf.Version
		Web struct {
			APIHOST         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
	}
	cfg.Version.Desc = "copyright information here"
	cfg.Version.SVN = BUILD

	if err := conf.Parse(os.Args[1:], NAMESPACE_CONF, &cfg); err != nil {
		switch err {
			case conf.ErrHelpWanted:
				usage, err := conf.Usage(NAMESPACE_CONF, &cfg)
				if err != nil {
					return errors.Wrap(err, "generating config usage")
				}
				fmt.Println(usage)
				return nil
			case conf.ErrVersionWanted:
				version, err := conf.VersionString(NAMESPACE_CONF, &cfg)
				if err != nil {
					return errors.Wrap(err, "generating config version")
				}
				fmt.Println(version)
				return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =================================================================================================================
	// App starting

	// Print the build version for our logs. ALso expose it under /debug/vars
	expvar.NewString("build").Set(BUILD)
	log.Printf("main: Started : Application initializing : version %q", BUILD)
	defer log.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config: \n%v\n", out)


	// =================================================================================================================
	// Start Debug Service
	//
	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	// /debug/vars - Added to the default mux by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdown.

	log.Println("main: Initializing debugging support.")

	go func() {
		log.Printf("main Debug listening %s", cfg.Web.DebugHost)
		if err := http.ListenAndServe(cfg.Web.DebugHost, http.DefaultServeMux); err != nil {
			log.Printf("main: Debug listener closed : %v", err)
		}
	}()

	// =================================================================================================================
	// Start API Service

	log.Println("main: Initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr: cfg.Web.APIHOST,
		Handler: nil,
		ReadTimeout: cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests
	go func() {
		log.Printf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =================================================================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
		case err := <-serverErrors:
			return errors.Wrap(err, "server error")

		case sig := <-shutdown:
			log.Printf("main: %v : Start shutdown", sig)

			// Give outstanding requests a deadline for completing.
			ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
			defer cancel()

			// Asking listener to shutdown and shed load.
			if err := api.Shutdown(ctx); err != nil {
				api.Close()
				return errors.Wrap(err, "could not stop server gracefully")
			}
	}

	return nil
}