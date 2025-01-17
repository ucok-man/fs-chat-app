include .env.local

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^//'

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# COMMAND
# ==================================================================================== #
# Default values for environment variables
PORT ?= 3000
ENV ?= development
LOG_LEVEL ?= INFO
DB_MAX_OPEN_CONN ?= 25
DB_MAX_IDLE_CONN ?= 25
DB_MAX_IDLE_TIME ?= 15m
TRUSTED_ORIGINS ?= ""
JWT_SECRET ?="secret"
MIGRATION_DIR ?= "migrations"

FLAGS=\
	--port=${PORT} \
	--env=${ENV} \
	--log-level=${LOG_LEVEL} \
	--db-dsn=${DB_DSN} \
	--db-max-open-conns=${DB_MAX_OPEN_CONN} \
	--db-max-idle-conns=${DB_MAX_IDLE_CONN} \
	--db-max-idle-time=${DB_MAX_IDLE_TIME} \
	--cors-trusted-origins=${TRUSTED_ORIGINS} \
	--jwt-secret=${JWT_SECRET}


## run/dev: compile and run application in watch mode
.PHONY: run/dev
run/dev:
	@go run ./cmd/api/... ${FLAGS}

## run/build: build the binary version of this application
.PHONY: run/build
run/build:
	@go build -o ./build/app ./cmd/api/...

## run/start: run the output binary
.PHONY: run/start
run/start:
	@./build/app ${FLAGS}


# ==================================================================================== #
# MIGRATION
# ==================================================================================== #

GOOSE_CFG=\
	GOOSE_DRIVER="postgres"\
	GOOSE_DBSTRING=${DB_DSN}\
	GOOSE_MIGRATION_DIR=${MIGRATION_DIR}

## migrate/new name=$1: create a new database migration
.PHONY: migrate/new
migrate/new:
	@${GOOSE_CFG} goose create ${name} sql

## migrate/up: run all up migration
.PHONY: migrate/up
migrate/up:
	@${GOOSE_CFG} goose up

## migrate/down: roll back all down migration
.PHONY: migrate/down
migrate/down:
	@${GOOSE_CFG} goose reset

## migrate/back: roll back migration by 1
.PHONY: migrate/back
migrate/back:
	@${GOOSE_CFG} goose down

## migrate/version: show current active migration
.PHONY: migrate/version
migrate/version:
	@${GOOSE_CFG} goose version


# ------------------------------------------------------------------ #
#                                 TLS                                #
# ------------------------------------------------------------------ #

## cert/linux: generate development certificate for linux os 
.PHONY: cert/linux
cert/linux:
	@go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
	@mv *.pem tls