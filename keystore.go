package main

import (
	"crypto/x509"
	"fmt"
	"os"
)

func loadKeystore(path string) (*x509.CertPool, error) {
	certs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	if certs == nil {
		certs = x509.NewCertPool()
	}

	caCert, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ok := certs.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf("failed to parse root certificate")
	}

	return certs, nil
}
