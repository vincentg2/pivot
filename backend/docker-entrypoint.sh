#!/bin/sh
set -eu

migration_database_url="${MIGRATION_DATABASE_URL:-$DATABASE_URL}"
goose -dir /app/migrations postgres "$migration_database_url" up
exec pivot-api
