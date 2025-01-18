package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/ucok-man/fs-chat-app-backend/internal/logger"
	"github.com/ucok-man/fs-chat-app-backend/internal/validator"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	log struct {
		level logger.Level
	}
	jwt struct {
		secret string
	}
	cloudinary struct {
		url string
	}
	cors struct {
		trustedOrigins []string
	}
}

// configuration - initialize and validate config.
// returning config and posible error validation in the form slice of key.
func configuration() (config, map[string]string) {
	var cfg config

	p := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	g := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	p.IntVar(&cfg.port, "port", 4000, "API server port")
	p.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production")

	p.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN <postgres://user:password@url/dbname?sslmode=disable>")
	p.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	p.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max open idle connections")
	defaultDuration, _ := time.ParseDuration("15m")
	p.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", defaultDuration, "PostgreSQL max connection idle time. Accept any valid parse duration value see https://pkg.go.dev/time#ParseDuration")

	cfg.log.level = logger.LevelInfo // default level
	g.Func("log-level", "Minimum log level (INFO|ERROR|FATAL) (default \"INFO\")", func(s string) error {
		switch s {
		case "INFO":
			cfg.log.level = logger.LevelInfo
		case "ERROR":
			cfg.log.level = logger.LevelError
		case "FATAL":
			cfg.log.level = logger.LevelFatal
		default:
			return fmt.Errorf("permitted value (INFO|ERROR|FATAL)")
		}
		return nil
	})

	p.StringVar(&cfg.jwt.secret, "jwt-secret", "", "JWT token secret")
	p.StringVar(&cfg.cloudinary.url, "cloudinary-url", "", "Cloudinary url <cloudinary://API_KEY:API_SECRET@CLOUD_NAME>")
	p.StringSliceVar(&cfg.cors.trustedOrigins, "cors-trusted-origins", []string{}, "Trusted CORS origins (comma separated)")

	p.AddGoFlagSet(g)
	p.Parse(os.Args[1:])

	v := validator.New()
	if validateConfig(v, &cfg); !v.Valid() {
		return config{}, v.Errors
	}

	return cfg, nil
}

// validateConfig - validate required config.
func validateConfig(v *validator.Validator, cfg *config) {
	v.Check(cfg.port > 0 && cfg.port <= 65535, "port", "invalid port number")
	v.Check(validator.PermittedValue(cfg.env, "development", "staging", "production"), "env", "invalid environment value")
	v.Check(cfg.db.dsn != "", "db-dsn", "databse source is required")
	v.Check(cfg.db.maxOpenConns > 0, "db-max-open-conns", "must be positive number")
	v.Check(cfg.db.maxIdleConns > 0, "db-max-idle-conns", "must be positive number")
	v.Check(validator.PermittedValue(cfg.log.level, logger.LevelError, logger.LevelFatal, logger.LevelInfo), "log-level", "invalid log level value")
	v.Check(cfg.jwt.secret != "", "jwt-secret", "jwt secret key is required")
	v.Check(cfg.cloudinary.url != "", "cloudinary-url", "cloudinary url is required")
}
