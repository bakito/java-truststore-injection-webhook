#!/bin/bash
set -e

kubectl apply -n java-truststore-injection-webhook -f testdata/e2e/cert-configmap.yaml

echo "Read cacerts"

CACERTS=$(kubectl get cm -n java-truststore-injection-webhook java-certs -o json | jq -r '.binaryData.cacerts')

echo "${CACERTS}" | base64 --decode > cacerts
keytool -list -keystore cacerts -storepass changeit | grep "Your keystore contains 1 entry"
