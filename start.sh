#!/bin/sh
DB_SOURCE=postgresql://nader:nader123@postgres:5432/billing_system?sslmode=disable
set -e

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"