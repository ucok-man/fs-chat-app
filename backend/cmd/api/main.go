package main

import (
	"errors"
	"os"
	"sync"

	"github.com/ucok-man/fs-chat-app-backend/internal/data"
	"github.com/ucok-man/fs-chat-app-backend/internal/logger"
)

const version = "1.0.0"

type application struct {
	config config
	logger *logger.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	l := logger.New(logger.WithLevel(logger.LevelInfo))
	l.SetDefault()

	cfg, errs := configuration()
	if errs != nil {
		l.Fatal(errors.New("invalid or missing flag")).Attr("meta", errs).Send()
		os.Exit(1)
	}

	db, err := opendb(cfg)
	if err != nil {
		l.Fatal(err).Attr("meta", "error opening database connection").Send()
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger.New(logger.WithLevel(cfg.log.level)),
		models: data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		l.Fatal(err).Attr("meta", "error running server").Send()
		os.Exit(1)
	}
}
