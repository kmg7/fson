package netutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path"
	"time"
)

func GenerateCert(out string, name string, org []string, notAfter time.Time) {
	// Generate a new private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Create a new certificate template
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   name,
			Organization: org,
		},
		NotBefore: time.Now(),
		NotAfter:  notAfter,
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
	}

	// Create a self-signed certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		panic(err)
	}

	// Save the certificate to a file
	certFile, err := os.Create(path.Join(out, "cert.pem"))
	if err != nil {
		panic(err)
	}
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certFile.Close()

	// Save the private key to a file
	keyFile, err := os.Create(path.Join(out, "key.pem"))
	if err != nil {
		panic(err)
	}
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	keyFile.Close()
}
