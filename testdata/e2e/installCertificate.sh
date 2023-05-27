#!/bin/bash
set -e

kubectl create ns java-truststore-injection-webhook || true
kubectl apply -n java-truststore-injection-webhook -f ./testdata/e2e/installCertificate.yaml
