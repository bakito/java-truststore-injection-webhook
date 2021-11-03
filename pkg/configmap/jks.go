package configmap

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/pavel-v-chernykh/keystore-go"
	"strings"
	"time"
)

func exportCerts(pems []*pem.Block, jksPassword string, t time.Time) ([]byte, error) {
	ks := keystore.KeyStore{}

	for i, p := range pems {

		ce := &keystore.TrustedCertificateEntry{
			Entry: keystore.Entry{
				CreationDate: t,
			},
			Certificate: keystore.Certificate{
				Content: p.Bytes,
				Type:    "X.509",
			},
		}
		ce.CreationDate = t

		ks[alias(p.Bytes, i)] = ce
	}

	var buf bytes.Buffer
	err := keystore.Encode(&buf, ks, []byte(jksPassword))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func alias(pem []byte, i int) string {
	c, err := x509.ParseCertificate(pem)
	if err != nil || c.Subject.CommonName == "" {
		return fmt.Sprintf("truststore-injector_%d", +i)
	}
	// inspired by: https://github.com/kaikramer/keystore-explorer/blob/79600e0e5cb5799dfc700df0989c5ba04f3d1db1/kse/src/org/kse/crypto/x509/X509CertUtil.java#L651

	if c.Issuer.CommonName == "" || c.Subject.CommonName == c.Issuer.CommonName {
		return strings.ToLower(c.Subject.CommonName)
	}
	return strings.ToLower(fmt.Sprintf("%s (%s)", c.Subject.CommonName, c.Issuer.CommonName))
}
