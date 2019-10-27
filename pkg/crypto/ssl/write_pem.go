package ssl

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

func NewPem(path string) error {
	pem, err := generatePem(KEY_LEN)
	if err != nil {
		return err
	}

	return writePem(path, pem)
}

func writePem(path string, buf *bytes.Buffer) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	return err
}

func generatePem(size int) (*bytes.Buffer, error) {
	key, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}

	cert, err := generateX509Cert(key)
	if err != nil {
		return nil, err
	}

	private_key, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "PRIVATE KEY", Bytes: private_key})
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: cert})

	return out, nil
}

func generateX509Cert(key *rsa.PrivateKey) ([]byte, error) {
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),

		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(EXPIRY),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	return x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
}
