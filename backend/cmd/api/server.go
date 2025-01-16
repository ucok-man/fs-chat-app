package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelDebug),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("caught signal").
			Attr("signal", s.String()).
			SkipC(2).
			Send()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks").
			Attr("addr", srv.Addr).
			SkipC(2).
			Send()

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.Info("starting server").
		Attr("addr", srv.Addr).
		Attr("env", app.config.env).
		SkipC(2).
		Send()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server").
		Attr("addr", srv.Addr).
		SkipC(2).
		Send()

	return nil
}
