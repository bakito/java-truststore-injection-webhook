# jti-webhook

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
    jtsi.bakito.ch/inject-truststore: 'true'
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