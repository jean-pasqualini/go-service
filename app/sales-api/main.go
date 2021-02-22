package main

import (
	"expvar"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
	"github.com/ardanlabs/conf"
)

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

	return nil
}