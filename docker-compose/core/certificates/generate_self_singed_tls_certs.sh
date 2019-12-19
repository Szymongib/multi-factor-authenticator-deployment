#!/usr/bin/env bash

set -e

CURRENT_DIR="$( cd "$(dirname "$0")" ; pwd -P )"

openssl genrsa -out tls.key 4096

openssl req -new -sha256 -key tls.key -out tls.csr -subj "/CN=localhost"

openssl req -x509 -sha256 -days 36500 -key tls.key -in tls.csr -out tls.crt
