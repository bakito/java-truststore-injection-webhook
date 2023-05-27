#!/bin/bash
set -e
sleep 10
kubectl apply -f testdata/e2e/cert-configmap.yaml
echo "Read cacerts"
CACERTS=$(kubectl get cm -n java-truststore-injection-webhook java-certs -o json | jq -r '.binaryData.cacerts')

echo "${CACERTS}" | base64 --decode > cacerts
keytool -list -v -keystore cacerts

