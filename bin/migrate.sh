#!/usr/bin/env sh

if [[ "$1" == "create" ]]; then
    exec migrate create -ext sql -dir database/migrations -seq "${@:2}"
else
    exec migrate
fi
