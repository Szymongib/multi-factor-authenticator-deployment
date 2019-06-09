#!/usr/bin/env bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"


kubectl apply -f ${SCRIPT_DIR}/tiller.yaml

${SCRIPT_DIR}/is-ready.sh kube-system name tiller


helm upgrade --install nginx-ingress stable/nginx-ingress --namespace default --set controller.replicaCount=1 --set controller.nodeSelector."beta\.kubernetes\.io/os"=linux --set defaultBackend.nodeSelector."beta\.kubernetes\.io/os"=linux --set controller.extraArgs.enable-ssl-passthrough=""

helm upgrade --install core ${SCRIPT_DIR}/../multi-factor-authenticator
