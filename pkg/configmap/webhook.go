package configmap

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/snorwin/k8s-generic-webhook/pkg/webhook"
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
	println(cm)
	return admission.Allowed("")
}
