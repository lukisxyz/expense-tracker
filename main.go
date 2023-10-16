package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/flukis/expt/service/internals/book"
	"github.com/flukis/expt/service/internals/config"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	var configFileName string
	flag.StringVar(
		&configFileName,
		"c",
		"config.yml",
		"Config file name",
	)
	flag.Parse()

	cfg := config.LoadConfig(configFileName)
	log.Debug().Any("config", cfg).Msg("config loaded")

	ctx := context.Background()

	pool, err := pgxpool.New(
		ctx,
		cfg.DBCfg.ConnStr(),
	)

	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
	}

	if err := book.SetPool(pool); err != nil {
		log.Error().Err(err).Msg("unable to set pool")
	}

	r := chi.NewRouter()
	r.Mount("/api/book", book.Router())

	log.Info().Msg(fmt.Sprintf("starting up server on: %s", cfg.Listen.Addr()))
	server := &http.Server{
		Handler:      r,
		Addr:         cfg.Listen.Addr(),
		ReadTimeout:  time.Second * time.Duration(cfg.Listen.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.Listen.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.Listen.IdleTimeout),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to start the server")
		return
	}
	log.Info().Msg("server stop")
}