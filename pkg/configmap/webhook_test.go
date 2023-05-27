package configmap_test

import (
	"context"
	"time"

	"github.com/bakito/java-truststore-injection-webhook/pkg/configmap"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			err := wh.Default(ctx, cm)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cm.BinaryData).Should(HaveLen(1))
			Ω(cm.BinaryData).Should(HaveKey(configmap.DefaultTruststoreName))
			// Ω(os.WriteFile("cacerts", cm.BinaryData["java-trust.jks"], 0644)).ShouldNot(HaveOccurred())
		})
		It("should add a cacerts binary entry with custom name", func() {
			cm.Data["a.pem"] = cert
			cm.Labels[configmap.LabelTruststoreName] = "java-trust.jks"
			err := wh.Default(ctx, cm)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cm.BinaryData).Should(HaveLen(1))
			Ω(cm.BinaryData).Should(HaveKey("java-trust.jks"))
		})
		It("should delete truststore with previous name", func() {
			cm.Data["a.pem"] = cert
			cm.BinaryData = map[string][]byte{"prev.jks": []byte("...")}
			cm.Annotations[configmap.AnnotationLastTruststoreName] = "prev.jks"
			err := wh.Default(ctx, cm)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cm.BinaryData).Should(HaveLen(1))
			Ω(cm.BinaryData).Should(HaveKey(configmap.DefaultTruststoreName))
		})

		It("should cacert must be reproducible", func() {
			cm.Data["a.pem"] = cert
			err := wh.Default(ctx, cm)
			Ω(err).ShouldNot(HaveOccurred())
			cacert1 := cm.BinaryData[configmap.DefaultTruststoreName]
			time.Sleep(3 * time.Second)
			err = wh.Default(ctx, cm)
			Ω(err).ShouldNot(HaveOccurred())
			cacert2 := cm.BinaryData[configmap.DefaultTruststoreName]
			Ω(cacert1).Should(Equal(cacert2))
		})

		It("should remove cacert if the label is missing", func() {
			delete(cm.Labels, configmap.LabelEnabled)
			cm.BinaryData = map[string][]byte{configmap.DefaultTruststoreName: []byte("test")}
			err := wh.Default(ctx, cm)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cm.BinaryData).Should(BeEmpty())
		})
	})
})

const (
	// openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=XX/ST=github/L=bakito/CN=java-truststore-injection-webhook"
	cert = `-----BEGIN CERTIFICATE-----
MIIFmTCCA4GgAwIBAgIUKvSUzRiN3GyPAJk+x7zywVJIjL8wDQYJKoZIhvcNAQEL
BQAwWzELMAkGA1UEBhMCWFgxDzANBgNVBAgMBmdpdGh1YjEPMA0GA1UEBwwGYmFr
aXRvMSowKAYDVQQDDCFqYXZhLXRydXN0c3RvcmUtaW5qZWN0aW9uLXdlYmhvb2sw
IBcNMjMwNTI3MTEyMDQ3WhgPMjEyMzA1MDMxMTIwNDdaMFsxCzAJBgNVBAYTAlhY
MQ8wDQYDVQQIDAZnaXRodWIxDzANBgNVBAcMBmJha2l0bzEqMCgGA1UEAwwhamF2
YS10cnVzdHN0b3JlLWluamVjdGlvbi13ZWJob29rMIICIjANBgkqhkiG9w0BAQEF
AAOCAg8AMIICCgKCAgEAq/7ZUh3aGkMTCzsoPbylqiBFsNNrPb/SBKpdljoKpejk
Kuc/OSjtjgIPSZODUxnW9p+vwJ5Sv8nbqOpwfgkFHUwItEC771NqOOfObheiZZGz
QxSqJUGActckTmGhgRC2TAWAFugzuusQK3EHAXOycBbDflTfC0IgEquUYExAQ2wB
fmAoImJPsUXpoKtqVL7BnTemPcl0eUp07jmsxWmwmTd5STF36UsRO5sydpr9JM1P
WltOfR+gNVoLU4A2mi8M5K+t0aBFgvZsvryzOVZG5m3RnUcd9guwDLTgg/hrXJdg
LPO5iubmq2haqNS++eDVmuGkZSM8V8y73BVPfwX1VwhrSyOk0qS6faM9jl3egDKx
layHjBpZEJzaAyGf8BiJqNCQa80qkRQ7Tkm2ZG8zHPSuvieHPncIQ45dcucJDHPe
E06V4Tk8vCbSOLyZeb0GtRCZej+TUc2u8OSaVaCHKDAuGYnBJwmdo8UWM2Ta19H3
Qp/7o+HTiiMi3TH68mb1M07hDuWmHpdhmE+4bG3/u0lFkis9PSyq9roiyuEjoZyB
5oZ+5LOQGbeXmwHG8F+9vHiR03PqxuS7mppOk399sTDN9Pr0RfjxLEH0lmOiGP2f
DJxA/dTU6xSK0CvFB+Ta9SwOHMX0bTBZMHJ0Ew4ATRAtBzkTahbOEk4h5p7eoyMC
AwEAAaNTMFEwHQYDVR0OBBYEFODsUAiQQz6fAD8+dgiloe7R6PocMB8GA1UdIwQY
MBaAFODsUAiQQz6fAD8+dgiloe7R6PocMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZI
hvcNAQELBQADggIBAE5KRdTrbXMMJUFAFmO2JPhoqTXSjA/rAe3gTcGTDy9uk1yD
caXRs22IhCofAjZzPZJv/efg98dWwhA3FPiHYwOg+efdqRIAjdUDfN4GOz70SOmb
fT0da4A+MOMDbGAjo1jz7HmXLxxzwSP1xVfQVNTklFuu6S4ymS+za8AApmflakRH
eIPOC9da1gfvjsE85a3PanzDYmSXEoo2wwULGI03ghQMq2VoO0Q7Kr7dW/HeumPX
ISbzsycPJSk0PDyaJbutz/wJE0Uv+rM81g3uWyfUCSfHTfganup5IkPmQcM4eFcV
GkOGs3RKfQul5rK/kG8VVyKqDVllZ6KUhodLZXTxMJGZEOlhNw3GIy50dbNhXJE5
myuQUUs5JBqfBKxBocfTy2mKe4WaGVeDPpodFnSFMdDdQQNBBwRmLmWJIwEKDpcW
FSaj7mKUe5RnvFZTeHMLHDm/LgfaJHM0wrQo9YcYLaVQ4iYVNE94iPkI2zuhOgRq
TYrDFyVqUA/jrgaAbguGiwP4mddfUmaWUiafLzdLxm+37V2bYT5b+Vn+PmPNpVom
sFjsEiqNgPMsb+ARrzXdJzIO3gEW0TYLAx24ABgRttd0TXURZbuGkVwoEyDSQHer
9tcIzkwuz9wxWrxygcKc5u3AlOSE5A2CBXEWEshgbegYc/8jhuYbtDgONvkO
-----END CERTIFICATE-----	`
)
