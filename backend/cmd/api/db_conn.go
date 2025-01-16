package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func opendb(cfg config) (*sql.DB, error) {
	dbconn, err := sql.Open("pgx", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	dbconn.SetMaxOpenConns(cfg.db.maxOpenConns)
	dbconn.SetMaxIdleConns(cfg.db.maxIdleConns)
	dbconn.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = dbconn.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return dbconn, nil
}
