package configmap_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bakito/java-truststore-injection-webhook/pkg/configmap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ = Describe("Configmap", func() {
	Context("Mutate", func() {
		var (
			ctx context.Context
			wh  *configmap.Webhook
			cm  *corev1.ConfigMap
		)
		BeforeEach(func() {
			ctx = context.TODO()
			wh = &configmap.Webhook{}
			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.Now(),
					Labels: map[string]string{
						configmap.LabelEnabled: "true",
					},
					Annotations: map[string]string{},
				},
				Data: make(map[string]string),
			}
		})
		It("should add a cacerts binary entry", func() {
			cm.Data["a.pem"] = cert
			// c, _ := ioutil.ReadFile("/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem")
			// cm.Data["a.pem"] = string(c)
			wh.Mutate(ctx, admission.Request{}, cm)
			Ω(cm.BinaryData).Should(HaveLen(1))
			Ω(cm.BinaryData).Should(HaveKey(configmap.DefaultTruststoreName))
			// Ω(os.WriteFile("cacerts", cm.BinaryData["java-trust.jks"], 0644)).ShouldNot(HaveOccurred())
		})
		It("should add a cacerts binary entry with custom name", func() {
			cm.Data["a.pem"] = cert
			cm.Labels[configmap.LabelTruststoreName] = "java-trust.jks"
			wh.Mutate(ctx, admission.Request{}, cm)
			Ω(cm.BinaryData).Should(HaveLen(1))
			Ω(cm.BinaryData).Should(HaveKey("java-trust.jks"))
		})
		It("should delete truststore with previous name", func() {
			cm.Data["a.pem"] = cert
			cm.BinaryData = map[string][]byte{"prev.jks": []byte("...")}
			cm.Annotations[configmap.AnnotationLastTruststoreName] = "prev.jks"
			wh.Mutate(ctx, admission.Request{}, cm)
			Ω(cm.BinaryData).Should(HaveLen(1))
			Ω(cm.BinaryData).Should(HaveKey(configmap.DefaultTruststoreName))
		})

		It("should cacert must be reproducable", func() {
			cm.Data["a.pem"] = cert
			wh.Mutate(ctx, admission.Request{}, cm)
			cacert1 := cm.BinaryData[configmap.DefaultTruststoreName]
			time.Sleep(3 * time.Second)
			wh.Mutate(ctx, admission.Request{}, cm)
			cacert2 := cm.BinaryData[configmap.DefaultTruststoreName]
			Ω(cacert1).Should(Equal(cacert2))
		})

		It("should remove cacert if the label is missing", func() {
			delete(cm.Labels, configmap.LabelEnabled)
			cm.BinaryData = map[string][]byte{configmap.DefaultTruststoreName: []byte("test")}
			wh.Mutate(ctx, admission.Request{}, cm)
			Ω(cm.BinaryData).Should(BeEmpty())
		})
	})
})

const (
	cert = `-----BEGIN CERTIFICATE-----
MIIH0zCCBbugAwIBAgIIXsO3pkN/pOAwDQYJKoZIhvcNAQEFBQAwQjESMBAGA1UE
AwwJQUNDVlJBSVoxMRAwDgYDVQQLDAdQS0lBQ0NWMQ0wCwYDVQQKDARBQ0NWMQsw
CQYDVQQGEwJFUzAeFw0xMTA1MDUwOTM3MzdaFw0zMDEyMzEwOTM3MzdaMEIxEjAQ
BgNVBAMMCUFDQ1ZSQUlaMTEQMA4GA1UECwwHUEtJQUNDVjENMAsGA1UECgwEQUND
VjELMAkGA1UEBhMCRVMwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCb
qau/YUqXry+XZpp0X9DZlv3P4uRm7x8fRzPCRKPfmt4ftVTdFXxpNRFvu8gMjmoY
HtiP2Ra8EEg2XPBjs5BaXCQ316PWywlxufEBcoSwfdtNgM3802/J+Nq2DoLSRYWo
G2ioPej0RGy9ocLLA76MPhMAhN9KSMDjIgro6TenGEyxCQ0jVn8ETdkXhBilyNpA
lHPrzg5XPAOBOp0KoVdDaaxXbXmQeOW1tDvYvEyNKKGno6e6Ak4l0Squ7a4DIrhr
IA8wKFSVf+DuzgpmndFALW4ir50awQUZ0m/A8p/4e7MCQvtQqR0tkw8jq8bBD5L/
0KIV9VMJcRz/RROE5iZe+OCIHAr8Fraocwa48GOEAqDGWuzndN9wrqODJerWx5eH
k6fGioozl2A3ED6XPm4pFdahD9GILBKfb6qkxkLrQaLjlUPTAYVtjrs78yM2x/47
4KElB0iryYl0/wiPgL/AlmXz7uxLaL2diMMxs0Dx6M/2OLuc5NF/1OVYm3z61PMO
m3WR5LpSLhl+0fXNWhn8ugb2+1KoS5kE3fj5tItQo05iifCHJPqDQsGH+tUtKSpa
cXpkatcnYGMN285J9Y0fkIkyF/hzQ7jSWpOGYdbhdQrqeWZ2iE9x6wQl1gpaepPl
uUsXQA+xtrn13k/c4LOsOxFwYIRKQ26ZIMApcQrAZQIDAQABo4ICyzCCAscwfQYI
KwYBBQUHAQEEcTBvMEwGCCsGAQUFBzAChkBodHRwOi8vd3d3LmFjY3YuZXMvZmls
ZWFkbWluL0FyY2hpdm9zL2NlcnRpZmljYWRvcy9yYWl6YWNjdjEuY3J0MB8GCCsG
AQUFBzABhhNodHRwOi8vb2NzcC5hY2N2LmVzMB0GA1UdDgQWBBTSh7Tj3zcnk1X2
VuqB5TbMjB4/vTAPBgNVHRMBAf8EBTADAQH/MB8GA1UdIwQYMBaAFNKHtOPfNyeT
VfZW6oHlNsyMHj+9MIIBcwYDVR0gBIIBajCCAWYwggFiBgRVHSAAMIIBWDCCASIG
CCsGAQUFBwICMIIBFB6CARAAQQB1AHQAbwByAGkAZABhAGQAIABkAGUAIABDAGUA
cgB0AGkAZgBpAGMAYQBjAGkA8wBuACAAUgBhAO0AegAgAGQAZQAgAGwAYQAgAEEA
QwBDAFYAIAAoAEEAZwBlAG4AYwBpAGEAIABkAGUAIABUAGUAYwBuAG8AbABvAGcA
7QBhACAAeQAgAEMAZQByAHQAaQBmAGkAYwBhAGMAaQDzAG4AIABFAGwAZQBjAHQA
cgDzAG4AaQBjAGEALAAgAEMASQBGACAAUQA0ADYAMAAxADEANQA2AEUAKQAuACAA
QwBQAFMAIABlAG4AIABoAHQAdABwADoALwAvAHcAdwB3AC4AYQBjAGMAdgAuAGUA
czAwBggrBgEFBQcCARYkaHR0cDovL3d3dy5hY2N2LmVzL2xlZ2lzbGFjaW9uX2Mu
aHRtMFUGA1UdHwROMEwwSqBIoEaGRGh0dHA6Ly93d3cuYWNjdi5lcy9maWxlYWRt
aW4vQXJjaGl2b3MvY2VydGlmaWNhZG9zL3JhaXphY2N2MV9kZXIuY3JsMA4GA1Ud
DwEB/wQEAwIBBjAXBgNVHREEEDAOgQxhY2N2QGFjY3YuZXMwDQYJKoZIhvcNAQEF
BQADggIBAJcxAp/n/UNnSEQU5CmH7UwoZtCPNdpNYbdKl02125DgBS4OxnnQ8pdp
D70ER9m+27Up2pvZrqmZ1dM8MJP1jaGo/AaNRPTKFpV8M9xii6g3+CfYCS0b78gU
JyCpZET/LtZ1qmxNYEAZSUNUY9rizLpm5U9EelvZaoErQNV/+QEnWCzI7UiRfD+m
AM/EKXMRNt6GGT6d7hmKG9Ww7Y49nCrADdg9ZuM8Db3VlFzi4qc1GwQA9j9ajepD
vV+JHanBsMyZ4k0ACtrJJ1vnE5Bc5PUzolVt3OAJTS+xJlsndQAJxGJ3KQhfnlms
tn6tn1QwIgPBHnFk/vk4CpYY3QIUrCPLBhwepH2NDd4nQeit2hW3sCPdK6jT2iWH
7ehVRE2I9DZ+hJp4rPcOVkkO1jMl1oRQQmwgEh0q1b688nCBpHBgvgW1m54ERL5h
I6zppSSMEYCUWqKiuUnSwdzRp+0xESyeGabu4VXhwOrPDYTkF7eifKXeVSUG7szA
h1xA2syVP1XgNce4hL60Xc16gwFy7ofmXx2utYXGJt/mwZrpHgJHnyqobalbz+xF
d3+YJ5oyXSrjhO7FmGYvliAd3djDJ9ew+f7Zfc3Qn48LFFhRny+Lwzgt3uiP1o2H
pPVWQxaZLPSkVrQ0uGE3ycJYgBugl6H8WY3pEfbRD0tVNEYqi4Y7
-----END CERTIFICATE-----`
)
