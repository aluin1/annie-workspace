#!/bin/bash

set -e

if [[ ! -d "/app/vendor" || -z $(ls -A "/app/vendor") ]]; then
	echo "[build.sh: go mod vendor for $APP_NAME]"
	cd /app && go mod vendor
fi

# Run compiled service
go run -mod=vendor /app/main.go "$@"


#!/bin/sh

set -e

echo "APP_NAME: $APP_NAME"

if [ ! -d "/app/vendor" ] || [ -z "$(ls -A /app/vendor)" ]; then
    echo "[build.sh: go mod vendor for $APP_NAME]"
    cd /app && go mod vendor
fi

if [ ! -f "/app/main.go" ]; then
    echo "Error: /app/main.go not found!"
    exit 1
fi

# Run compiled service
go run -mod=vendor /app/main.go "$@"
