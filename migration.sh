#!/bin/sh
set -e
export MIGRATION_DIR=./migrations

if [ -z "$PG_DSN" ]; then
    export DB_NAME="${DB_NAME:-notes_db}"
    export DB_HOST="${DB_HOST:-localhost}"
    export DB_PORT="${DB_PORT:-5432}"
    export DB_USER="${DB_USER:-notes_user}"
    export DB_PASSWORD="${DB_PASSWORD:-notes_password}"
    export DB_SSL="${DB_SSLMODE:-disable}"

    export PG_DSN="host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} user=${DB_USER} sslmode=${DB_SSL} password=${DB_PASSWORD}"
fi

if [ "$1" = '--dryrun' ]; then
    ./bin/goose -allow-missing -dir $MIGRATION_DIR postgres "${PG_DSN}" status -v
else
    ./bin/goose -allow-missing -dir $MIGRATION_DIR postgres "${PG_DSN}" up -s
fi