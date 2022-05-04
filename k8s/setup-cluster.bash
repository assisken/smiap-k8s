#!/usr/bin/env bash

set -o pipefail -o xtrace -o errexit

CLUSTER_NAME="smiap"

SCRIPT_PATH="$(realpath .)"

# Create registry container unless it already exists
if [ "$(docker inspect -f '{{.State.Running}}' "kind-registry" 2>/dev/null || true)" != 'true' ]; then
  docker run -d --restart=always -p "127.0.0.1:5000:5000" --name "kind-registry" registry:2
fi

kind get clusters | grep "$CLUSTER_NAME" || kind create cluster --config "$SCRIPT_PATH/k8s/config/cluster.yaml" --wait 60s
rm -f .kubeconfig

# Connect the registry to the cluster network if not already connected
docker network connect "kind" "kind-registry" || true

kubeconfig="$(mktemp)"
trap 'rm -f "$kubeconfig"' EXIT
kind get kubeconfig --name smiap > "$kubeconfig"

export KUBECONFIG="$kubeconfig"
