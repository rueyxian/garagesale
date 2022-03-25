package main

import (
	"context"
	_ "expvar" // register the /debug/vars handler
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // register the /debug/pprof handler
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/naixyeur/garagesale/cmd/sales-api/internal/handlers"
	"github.com/naixyeur/garagesale/internal/platform/conf"
	"github.com/naixyeur/garagesale/internal/platform/database"
	"github.com/pkg/errors"
)

// ================================================================================

func main() {

	if err := run(); err != nil {
		log.Fatal(err)
	}

}

func run() error {
	log := log.New(os.Stdout, "SALES : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	log.Printf("started\n")
	defer log.Printf("completed\n")

	// ==============================
	// environment variables
	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			Debug           string        `conf:"default:localhost:6060"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:true"`
		}
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// ==============================
	// debug service
	go func() {
		log.Printf("debug server listening on %s\n", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
		log.Printf("debug server close %s", err)
	}()

	// ==============================
	// api service
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer db.Close()

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handlers.API(log, db),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverError := make(chan error, 1)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// ==============================

	go func() {
		log.Printf("api server listening on %s\n", api.Addr)
		serverError <- api.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		return errors.Wrap(err, "listening and serving")
	case <-shutdown:
		fmt.Println()
		log.Printf("shutdown in %v\n", cfg.Web.ShutdownTimeout)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err := api.Shutdown(ctx)

		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}

		if err != nil {
			return errors.Wrap(err, "graceful shutdown")
		}
	}
	return nil
}

// ================================================================================
