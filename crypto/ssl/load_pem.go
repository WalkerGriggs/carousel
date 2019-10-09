package ssl

import (
	"crypto"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func LoadPem(path string) (*tls.Certificate, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cert tls.Certificate
	for {
		p, rest := pem.Decode(raw)
		if p == nil {
			break
		}

		switch p.Type {
		case "CERTIFICATE":
			cert.Certificate = append(cert.Certificate, p.Bytes)

		case "PRIVATE KEY", "RSA PRIVATE KEY":
			cert.PrivateKey, err = parsePrivateKey(p.Bytes)
			if err != nil {
				return nil, err
			}
		}
		raw = rest
	}

	return &cert, nil
}

func parsePrivateKey(raw []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS8PrivateKey(raw); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey:
			return key, nil

		default:
			return nil, fmt.Errorf("Unkown or unsupported private key type in PKCS8 wrapper")
		}
	}

	return nil, fmt.Errorf("Unsupported private key wrapper")
}
