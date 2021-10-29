package configmap

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"github.com/pavel-v-chernykh/keystore-go"
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

		ks[fmt.Sprintf("truststore-injector_%d", +i)] = ce
	}

	var buf bytes.Buffer
	err := keystore.Encode(&buf, ks, []byte(jksPassword))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
