#!/bin/sh

set -e

echo "run db migration"
migrate -database postgres://admin:secret@postgres:5432/blog_db?sslmode=disable -path migrations up

echo "start the app"
exec "$@"