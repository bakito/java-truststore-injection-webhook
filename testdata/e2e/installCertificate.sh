#!/bin/bash
set -e

kubectl create namespace java-truststore-injection-webhook || true
kubectl apply -f ./testdata/e2e/certificate.yaml
