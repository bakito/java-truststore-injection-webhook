#!/bin/bash
set -e

kubectl create ns java-truststore-injection-webhook || true
kubectl apply -f ./testdata/e2e/certificate.yaml
