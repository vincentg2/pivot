#!/bin/sh
set -eu

goose -dir /app/migrations postgres "$DATABASE_URL" up
exec pivot-api
