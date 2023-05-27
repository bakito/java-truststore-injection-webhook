#!/bin/bash
set -e

helm upgrade --install java-truststore-injection-webhook charts/java-truststore-injection-webhook \
  --namespace java-truststore-injection-webhook \
  --create-namespace \
  -f testdata/e2e/e2e-values.yaml \
  --atomic
