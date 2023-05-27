#!/bin/bash
set -e

kubectl apply -n java-truststore-injection-webhook -f testdata/e2e/testConfigMaps.yaml

textConfigMap () {
  echo "Read ${2} from ${1}"
  CACERTS=$(kubectl get cm -n java-truststore-injection-webhook ${1} -o json | jq -r ".binaryData.\"${2}\"")
  echo "${CACERTS}" | base64 --decode > cacerts
  keytool -list -keystore cacerts -storepass ${3} | grep "Your keystore contains 1 entry"

  echo "Check last-injected-truststore-name"
  kubectl get cm -n java-truststore-injection-webhook ${1} -o json | jq -r '.metadata.annotations."jti.bakito.ch/last-injected-truststore-name"' | grep "${2}"
}

textConfigMap java-certs-simple cacerts changeit

textConfigMap java-certs-extended my-certs my-precious
