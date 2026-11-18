#/usr/bin/env sh

jq -s '.[0] * .[1]' gokrazy/base.json gokrazy/secrets.json > gokrazy/config.json
exec gok -i gokrazy --parent_dir . "$@"
