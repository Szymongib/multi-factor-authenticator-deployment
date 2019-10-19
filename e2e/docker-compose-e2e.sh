#!/usr/bin/env bash

echo "---------------WARNING---------------"
echo "This script might not work correctly"
echo "---------------WARNING---------------"

set -e

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

pushd ${CURRENT_DIR}/../docker-compose

./run.sh &> ${CURRENT_DIR}/test.log &
COMPOSE_PID=$!

trap "echo stoping; kill ${COMPOSE_PID}" EXIT

sleep 10

popd

pushd ${CURRENT_DIR}/test

echo "Running tests..."
GO11MODULE=on go test -v ./...

exit_code=$?

popd

exit ${exit_code}
