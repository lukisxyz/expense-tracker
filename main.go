package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/flukis/expt/service/internals/account"
	"github.com/flukis/expt/service/internals/book"
	"github.com/flukis/expt/service/internals/category"
	"github.com/flukis/expt/service/internals/config"
	customMiddleware "github.com/flukis/expt/service/internals/middleware"
	"github.com/flukis/expt/service/internals/record"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	account.SetTokenizeConfig(
		int(cfg.TokenCfg.AccessDuration),
		int(cfg.TokenCfg.RefreshDuration),
		cfg.TokenCfg.Key,
	)
	customMiddleware.SetTokenizeConfig(cfg.TokenCfg.Key)

	ctx := context.Background()

	pool, err := pgxpool.New(
		ctx,
		cfg.DBCfg.ConnStr(),
	)

	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
	}

	if err := book.SetPool(pool); err != nil {
		log.Error().Err(err).Msg("unable to set pool on book")
	}
	if err := category.SetPool(pool); err != nil {
		log.Error().Err(err).Msg("unable to set pool on category")
	}
	if err := record.SetPool(pool); err != nil {
		log.Error().Err(err).Msg("unable to set pool on record")
	}
	if err := account.SetPool(pool); err != nil {
		log.Error().Err(err).Msg("unable to set pool on account")
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Mount("/api/account", account.Router())
	r.Group(func(r chi.Router) {
		r.Use(customMiddleware.AuthJwt)
		r.Mount("/api/book", book.Router())
		r.Mount("/api/category", category.Router())
		r.Mount("/api/record", record.Router())
	})

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
