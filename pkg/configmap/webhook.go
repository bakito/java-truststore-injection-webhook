package configmap

import (
	"context"
	"encoding/pem"
	"strings"

	"github.com/bakito/cacert-truststore-webhook/pkg/jks"
	"github.com/snorwin/k8s-generic-webhook/pkg/webhook"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	DefaultTruststoreName = "cacerts"

	LabelEnabled = "truststore.bakito.ch/enabled"

	annotationTruststoreName = "truststore.bakito.ch/fileName"
	annotationTruststorePass = "truststore.bakito.ch/password"
)

type Webhook struct {
	webhook.MutatingWebhook
}

func (w *Webhook) SetupWebhookWithManager(mgr manager.Manager) error {
	return webhook.NewGenericWebhookManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		WithMutatePath("/mutate").
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

	// delete if the label is not present anymore
	if !isEnabled(cm) {
		delete(cm.BinaryData, tsn)
		return admission.Allowed("")
	}

	var allPems []*pem.Block
	for name, content := range cm.Data {
		if strings.HasSuffix(name, ".pem") {
			allPems = append(allPems, readCerts(content)...)
		}
	}

	b, _ := jks.ExportCerts(allPems, pass, cm.ObjectMeta.CreationTimestamp.Time)

	if cm.BinaryData == nil {
		cm.BinaryData = make(map[string][]byte)
	}
	cm.BinaryData[tsn] = b
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

func isEnabled(cm *corev1.ConfigMap) bool {
	if cm.Labels == nil {
		return false
	}
	_, ok := cm.Labels[LabelEnabled]
	return ok
}
