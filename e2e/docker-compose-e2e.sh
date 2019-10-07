#!/usr/bin/env bash

set -e

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

pushd ${CURRENT_DIR}/../docker-compose

./run.sh &
COMPOSE_PID=$!

sleep 15

popd

pushd ${CURRENT_DIR}/test

go clean --testcache

go test -v ./...

popd

kill ${COMPOSE_PID}