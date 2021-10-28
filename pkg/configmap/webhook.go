package configmap

import (
	"context"
	"encoding/pem"
	"fmt"
	"github.com/bakito/cacert-truststore-webhook/pkg/jks"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/snorwin/k8s-generic-webhook/pkg/webhook"
)

const (
	annotationTruststoreName = "ch.bakito/truststore/fileName"
	annotationTruststorePass = "ch.bakito/truststore/password"
)

type Webhook struct {
	webhook.MutatingWebhook
}

func (w *Webhook) SetupWebhookWithManager(mgr manager.Manager) error {
	return webhook.NewGenericWebhookManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Complete(w)
}

func (w *Webhook) Mutate(ctx context.Context, _ admission.Request, object runtime.Object) admission.Response {
	_ = log.FromContext(ctx)

	cm := object.(*corev1.ConfigMap)

	pass := "changeit"
	if p, ok := cm.Annotations[annotationTruststorePass]; ok {
		pass = p
	}
	fileName := "cacerts"
	if n, ok := cm.Annotations[annotationTruststoreName]; ok {
		fileName = n
	}

	var allPems []*pem.Block
	// delete all errors
	for name := range cm.Data {
		if strings.HasSuffix(name, ".error") {
			delete(cm.Data, name)
		}
	}
	for name, content := range cm.Data {
		if strings.HasSuffix(name, ".pem") {
			pems, err := readCerts(content)
			if err == nil {
				allPems = append(allPems, pems...)
			} else {
				cm.Data[fmt.Sprintf("%s.error", name)] = err.Error()
			}
		}
	}

	b, _ := jks.ExportCerts(allPems, pass)

	if cm.BinaryData == nil {
		cm.BinaryData = make(map[string][]byte)
	}
	cm.BinaryData[fileName] = b
	return admission.Allowed("")
}

func readCerts(certFile string) ([]*pem.Block, error) {
	raw := []byte(certFile)
	var pems []*pem.Block
	for {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			pems = append(pems, block)
		}
		raw = rest
	}

	return pems, nil
}
