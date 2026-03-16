package main

import (
	"crypto/x509"
	"fmt"
	"os"
)

// loadKeystore loads a PEM-encoded certificate file from the specified path
// and appends it to the system's root certificate pool.
func loadKeystore(path string) (*x509.CertPool, error) {
	// Initialize the certificate pool with system-trusted certificates
	certs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	// Fallback to an empty pool if system certificates cannot be loaded
	if certs == nil {
		certs = x509.NewCertPool()
	}

	// Read the certificate file from the provided path.
	// #nosec G304: Path is user-provided for testing purposes.
	caCert, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse and append the certificates to the pool
	ok := certs.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf("failed to parse root certificate")
	}

	return certs, nil
}
