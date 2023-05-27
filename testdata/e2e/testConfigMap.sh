#!/bin/bash
set -e

kubectl apply -f testdata/e2e/cert-configmap.yaml
echo "Read cacerts"
kubectl get cm java-certs -o json | jq -r '.binaryData.cacerts'

