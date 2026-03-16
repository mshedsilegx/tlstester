package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// CreateTLSConfig initializes a tls.Config based on the provided application configuration.
func CreateTLSConfig(cfg *Config) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		// #nosec G402: InsecureSkipVerify is configurable for testing purposes.
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}

	// Set ServerName for SNI if verification is enabled
	if !cfg.InsecureSkipVerify {
		hostname, _, err := net.SplitHostPort(cfg.Hostport)
		if err != nil {
			return nil, fmt.Errorf("error parsing hostport: %w", err)
		}
		tlsConfig.ServerName = hostname
	}

	// Configure specific TLS version if requested
	if cfg.TLSVersion != "" {
		switch cfg.TLSVersion {
		case "TLS1.0":
			tlsConfig.MinVersion = tls.VersionTLS10
			tlsConfig.MaxVersion = tls.VersionTLS10
		case "TLS1.1":
			tlsConfig.MinVersion = tls.VersionTLS11
			tlsConfig.MaxVersion = tls.VersionTLS11
		case "TLS1.2":
			tlsConfig.MinVersion = tls.VersionTLS12
			tlsConfig.MaxVersion = tls.VersionTLS12
		case "TLS1.3":
			tlsConfig.MinVersion = tls.VersionTLS13
			tlsConfig.MaxVersion = tls.VersionTLS13
		default:
			return nil, fmt.Errorf("invalid TLS version: %s", cfg.TLSVersion)
		}
	}

	// Configure specific cipher suite if requested
	if cfg.CipherSuite != "" {
		var ciphers []uint16
		for _, suite := range tls.CipherSuites() {
			if suite.Name == cfg.CipherSuite {
				ciphers = append(ciphers, suite.ID)
				break
			}
		}
		if len(ciphers) == 0 {
			return nil, fmt.Errorf("invalid Cipher Suite: %s", cfg.CipherSuite)
		}
		tlsConfig.CipherSuites = ciphers
	}

	// Load custom CA certificates if a keystore path is provided
	if cfg.Keystore != "" {
		certs, err := loadKeystore(cfg.Keystore)
		if err != nil {
			return nil, fmt.Errorf("error loading keystore: %w", err)
		}
		tlsConfig.RootCAs = certs
	}

	return tlsConfig, nil
}

// ConnectTLS establishes a TCP connection and performs the TLS handshake with retries.
func ConnectTLS(cfg *Config, tlsConfig *tls.Config) (*tls.Conn, error) {
	var conn net.Conn
	var err error

	// Attempt to connect with retries
	for i := 0; i <= cfg.Retries; i++ {
		conn, err = net.DialTimeout("tcp", cfg.Hostport, cfg.Timeout)
		if err == nil {
			break
		}
		if i < cfg.Retries {
			fmt.Printf("Connection failed (attempt %d): %v, retrying...\n", i+1, err)
			time.Sleep(1 * time.Second)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("connection failed after %d retries: %w", cfg.Retries+1, err)
	}

	// Wrap the TCP connection with TLS and perform handshake
	tlsConn := tls.Client(conn, tlsConfig)
	err = tlsConn.Handshake()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("TLS Handshake failed: %w", err)
	}

	return tlsConn, nil
}

// PrintConnectionState displays details about the established TLS connection.
func PrintConnectionState(conn *tls.Conn) {
	cert := conn.ConnectionState().PeerCertificates[0]
	fmt.Println("Certificate Information:")
	fmt.Println("  Subject:", cert.Subject)
	fmt.Println("  SAN:", cert.DNSNames)
	fmt.Println("  Issuer:", cert.Issuer)

	// Display Root CA if present in the chain
	if len(conn.ConnectionState().PeerCertificates) > 1 {
		fmt.Println("  Root CA:", conn.ConnectionState().PeerCertificates[len(conn.ConnectionState().PeerCertificates)-1].Subject)
	}

	fmt.Println("  Expiration Date:", cert.NotAfter.Format(time.RFC3339))
	fmt.Println("  Protocol:", tls.VersionName(conn.ConnectionState().Version))
	fmt.Println("  Cipher Suite:", tls.CipherSuiteName(conn.ConnectionState().CipherSuite))
}
