package jks

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pavel-v-chernykh/keystore-go"
)

var (
	out io.Writer = os.Stdout // modified during testing
)

func ExportCerts(pems []*pem.Block, jksPassword string, t time.Time) ([]byte, error) {
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

		ks[fmt.Sprintf("cert_%d", +i)] = ce
	}

	var buf bytes.Buffer
	err := keystore.Encode(&buf, ks, []byte(jksPassword))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func alias(cert *x509.Certificate) string {
	return fmt.Sprintf("%s (%s)", strings.ToLower(cert.Subject.CommonName), strings.ToLower(cert.Issuer.CommonName))
}

func closeIt(s *os.File) {
	_ = s.Close()
}
