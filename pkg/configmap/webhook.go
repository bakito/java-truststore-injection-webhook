package configmap

import (
	"context"
	"encoding/pem"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	// DefaultTruststoreName default name of the generated truststore
	DefaultTruststoreName = "cacerts"

	// LabelEnabled k8s label to enable java truststore injection
	LabelEnabled = "jti.bakito.ch/inject-truststore"
	// LabelTruststoreName k8s label to define a custom truststore name
	LabelTruststoreName = "jti.bakito.ch/truststore-name"

	annotationTruststorePass = "jti.bakito.ch/truststore-password" // #nosec G101
	// AnnotationLastTruststoreName k8s annotation to keep track of the last truststore name
	AnnotationLastTruststoreName = "jti.bakito.ch/last-injected-truststore-name"
)

var certsInConfigMap = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "jti_certificates_truststore",
		Help: "Number certificates in the truststore",
	},
	[]string{"namespace", "configmap", "truststore"},
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(certsInConfigMap)
}

// Webhook implementation
type Webhook struct {
	admission.CustomDefaulter
}

// SetupWebhookWithManager setup this webhook
func (w *Webhook) SetupWebhookWithManager(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		WithDefaulter(w).
		Complete()
}

// Default mutate the configmap
func (w *Webhook) Default(ctx context.Context, obj runtime.Object) error {
	cm := obj.(*corev1.ConfigMap)

	l := log.FromContext(ctx).WithValues("configmap", cm.Name)
	// req, err := admission.RequestFromContext(ctx);

	tsn := DefaultTruststoreName
	if cm.Labels != nil {
		if n, ok := cm.Labels[LabelTruststoreName]; ok {
			tsn = n
		}
	}

	pass := "changeit"
	if cm.Annotations != nil {
		if p, ok := cm.Annotations[annotationTruststorePass]; ok {
			pass = p
		}
		if ltn, ok := cm.Annotations[AnnotationLastTruststoreName]; ok && cm.BinaryData != nil {
			delete(cm.BinaryData, ltn)
			certsInConfigMap.DeleteLabelValues(cm.Namespace, cm.Name, ltn)
		}
	}

	// delete if the label is not present anymore
	if !isEnabled(cm) {
		l.Info("removing truststore")
		if cm.BinaryData != nil {
			delete(cm.BinaryData, tsn)
		}
		if cm.Annotations != nil {
			delete(cm.Annotations, AnnotationLastTruststoreName)
		}
		certsInConfigMap.DeleteLabelValues(cm.Namespace, cm.Name, tsn)
		return nil
	}

	var allPems []*pem.Block
	for name, content := range cm.Data {
		pems := readCerts(content)
		l.WithValues("fileName", name, "certs", len(pems)).V(3).Info("found certs")
		allPems = append(allPems, pems...)
	}

	b, _ := exportCerts(allPems, pass, cm.CreationTimestamp.Time)

	if cm.BinaryData == nil {
		cm.BinaryData = make(map[string][]byte)
	}
	cm.BinaryData[tsn] = b
	if cm.Annotations == nil {
		cm.Annotations = make(map[string]string)
	}
	cm.Annotations[AnnotationLastTruststoreName] = tsn
	l.WithValues("certs", len(allPems), "truststore", tsn).Info("added certs to truststore")
	certsInConfigMap.WithLabelValues(cm.Namespace, cm.Name, tsn).Set(float64(len(allPems)))

	return nil
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
	value, ok := cm.Labels[LabelEnabled]
	return ok && strings.EqualFold("true", value)
}
