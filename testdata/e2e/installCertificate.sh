#!/bin/bash
set -e

if ! kubectl get ns jti > /dev/null; then
  kubectl create ns jti
fi
kubectl apply -n jti -f ./testdata/e2e/installCertificate.yaml
