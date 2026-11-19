#/usr/bin/env sh

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
PROJECT_DIR="$SCRIPT_DIR/.."

compile=0
for item in "$@"; do
    if [[ "$item" == "overwrite" || "$item" == "update" ]]; then
        compile=1
    fi
done

if [[ "$compile" == "1" ]]; then
    set -x
    # replace secrets and use real version
    jq -s '.[0] * .[1]' $PROJECT_DIR/config.json $PROJECT_DIR/secrets.json > $PROJECT_DIR/gokrazy/config.json
    cp $PROJECT_DIR/static-dhcp-lease.json $PROJECT_DIR/wifi.json $PROJECT_DIR/mediamtx.yml $PROJECT_DIR/go2rtc.yaml  $PROJECT_DIR/gokrazy/
    exec gok -i gokrazy --parent_dir $PROJECT_DIR "$@"
else
    # use project config
    exec gok -i "inshrimpentation" --parent_dir "$PROJECT_DIR/.." "$@"
fi
