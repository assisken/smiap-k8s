#!/usr/bin/env bash

set -o pipefail -o xtrace -o errexit

CLUSTER_NAME="smiap"

SCRIPT_PATH="$(realpath .)"
CONTEXT="kind-$CLUSTER_NAME"

SMIAP_K8S_TAG="smiap-k8s"
docker build --tag "localhost:5000/$SMIAP_K8S_TAG" --tag "$SMIAP_K8S_TAG" .
docker push "localhost:5000/$SMIAP_K8S_TAG"

kubectl --context "$CONTEXT" apply -f "$SCRIPT_PATH/k8s/manifests"
