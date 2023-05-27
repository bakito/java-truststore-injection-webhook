#!/bin/bash
set -e

helm upgrade --install java-truststore-injection-webhook chart \
  --namespace jti \
  --create-namespace \
  -f testdata/e2e/e2e-values.yaml \
  --atomic
kubectl get pods -n jti
