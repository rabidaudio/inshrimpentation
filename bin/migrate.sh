#!/usr/bin/env sh

exec migrate -ext sql -dir database/migrations "$@"
