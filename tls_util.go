package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

func CreateTLSConfig(cfg *Config) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}

	if !cfg.InsecureSkipVerify {
		hostname, _, err := net.SplitHostPort(cfg.Hostport)
		if err != nil {
			return nil, fmt.Errorf("error parsing hostport: %w", err)
		}
		tlsConfig.ServerName = hostname
	}

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

	if cfg.Keystore != "" {
		certs, err := loadKeystore(cfg.Keystore)
		if err != nil {
			return nil, fmt.Errorf("error loading keystore: %w", err)
		}
		tlsConfig.RootCAs = certs
	}

	return tlsConfig, nil
}

func ConnectTLS(cfg *Config, tlsConfig *tls.Config) (*tls.Conn, error) {
	var conn net.Conn
	var err error

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

	tlsConn := tls.Client(conn, tlsConfig)
	err = tlsConn.Handshake()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("TLS Handshake failed: %w", err)
	}

	return tlsConn, nil
}

func PrintConnectionState(conn *tls.Conn) {
	cert := conn.ConnectionState().PeerCertificates[0]
	fmt.Println("Certificate Information:")
	fmt.Println("  Subject:", cert.Subject)
	fmt.Println("  SAN:", cert.DNSNames)
	fmt.Println("  Issuer:", cert.Issuer)

	if len(conn.ConnectionState().PeerCertificates) > 1 {
		fmt.Println("  Root CA:", conn.ConnectionState().PeerCertificates[len(conn.ConnectionState().PeerCertificates)-1].Subject)
	}

	fmt.Println("  Expiration Date:", cert.NotAfter.Format(time.RFC3339))
	fmt.Println("  Protocol:", tls.VersionName(conn.ConnectionState().Version))
	fmt.Println("  Cipher Suite:", tls.CipherSuiteName(conn.ConnectionState().CipherSuite))
}
