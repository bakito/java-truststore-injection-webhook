module github.com/bakito/cacert-truststore-webhook

go 1.16

require (
	github.com/bakito/cert-fetcher v1.0.1
	github.com/pavel-v-chernykh/keystore-go v2.1.0+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/snorwin/k8s-generic-webhook v1.2.4
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/controller-runtime v0.10.2
)
