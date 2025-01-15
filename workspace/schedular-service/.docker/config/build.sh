#!/bin/bash

set -e

if [[ ! -d "/app/vendor" || -z $(ls -A "/app/vendor") ]]; then
	echo "[build.sh: go mod vendor for $APP_NAME]"
	cd /app && go mod vendor
fi

# Run compiled service
go run -mod=vendor /app/main.go "$@"
