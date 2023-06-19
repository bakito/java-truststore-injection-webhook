[![Go Report Card](https://goreportcard.com/badge/github.com/bakito/java-truststore-injection-webhook)](https://goreportcard.com/report/github.com/bakito/java-truststore-injection-webhook)
[![Github Build](https://github.com/bakito/java-truststore-injection-webhook/actions/workflows/build.yml/badge.svg)](https://github.com/bakito/java-truststore-injection-webhook/actions/workflows/build.yml)
[![GitHub Release](https://img.shields.io/github/release/bakito/java-truststore-injection-webhook.svg?style=flat)](https://github.com/bakito/java-truststore-injection-webhook/releases)

# Java Truststore Injection Webhook

This webhook injects a java truststore into a k8s ConfigMap containing pem certificates. If a ConfigMap is labelled to
be injected with a java truststore, the webhook checks all existing data entries for pem certificates and adds all found
fount to a java truststore file that is added as binary data.

## Usage

Label a configmap where a java truststore should be injected.

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  labels:
    jti.bakito.ch/inject-truststore: 'true'
```

## truststore file name

The default truststore file name is '__cacerts__'

A different ConfigMap file name can be defined by adding the following __label__.

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  labels:
    jti.bakito.ch/truststore-name: 'custom-truststore-name'
```

## truststore password

The default truststore password is '__changeit__'

A different ConfigMap file name can be defined by adding the following __annotation__.

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  annotations:
    jti.bakito.ch/truststore-password": 'custom-password'
```

## Installation

**java-truststore-injection-webhook** can be installed via our Helm chart:

```sh
helm repo add bakito https://bakito.github.io/helm-charts
helm repo update

helm upgrade --install java-truststore-injection-webhook bakito/java-truststore-injection-webhook
```
