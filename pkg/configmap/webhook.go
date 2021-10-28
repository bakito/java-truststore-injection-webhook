package configmap

import (
	"context"
	"encoding/pem"
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
	DefaultTruststoreName = "cacerts"

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
	tsn := DefaultTruststoreName
	if n, ok := cm.Annotations[annotationTruststoreName]; ok {
		tsn = n
	}

	var allPems []*pem.Block
	for name, content := range cm.Data {
		if strings.HasSuffix(name, ".pem") {
			allPems = append(allPems, readCerts(content)...)
		}
	}

	if len(allPems) > 0 {
		b, _ := jks.ExportCerts(allPems, pass)

		if cm.BinaryData == nil {
			cm.BinaryData = make(map[string][]byte)
		}
		cm.BinaryData[tsn] = b
	} else {
		delete(cm.BinaryData, tsn)
	}
	return admission.Allowed("")
}

func readCerts(certFile string) []*pem.Block {
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

	return pems
}
