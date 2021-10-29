module github.com/bakito/truststore-injector-webhook

go 1.17

require (
	cloud.google.com/go v0.81.0 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/pavel-v-chernykh/keystore-go v2.1.0+incompatible
	github.com/snorwin/k8s-generic-webhook v1.2.4
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602 // indirect
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/controller-runtime v0.10.2
)
