package main

import (
	"flag"
	"fmt"
	"time"
)

// Config holds the application configuration parameters for the TLS test.
type Config struct {
	Hostport           string        // Target address in host:port format
	Timeout            time.Duration // Connection timeout
	Retries            int           // Number of connection attempts
	TLSVersion         string        // Requested TLS version (e.g., "TLS1.2")
	CipherSuite        string        // Specific cipher suite name
	Keystore           string        // Path to custom CA certificate file
	InsecureSkipVerify bool          // Whether to skip server certificate validation
	Version            bool          // Whether to just print the version and exit
}

// ParseFlags parses command-line flags and returns a Config pointer.
// It also handles the basic validation of required flags.
func ParseFlags(version string) (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Hostport, "hostport", "", "Hostname:port to test (required)")
	flag.DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "Timeout for connection")
	flag.IntVar(&cfg.Retries, "retries", 3, "Number of retries")
	flag.StringVar(&cfg.TLSVersion, "tls", "", "Specific TLS version (e.g., TLS1.2, TLS1.3)")
	flag.StringVar(&cfg.CipherSuite, "cipher", "", "Specific Cipher Suite (e.g., TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256)")
	flag.StringVar(&cfg.Keystore, "keystore", "", "Path to trusted keystore (optional)")
	flag.BoolVar(&cfg.InsecureSkipVerify, "insecure", false, "Skip chain of trust verification (use with caution)")
	flag.BoolVar(&cfg.Version, "version", false, "Display version information")

	flag.Parse()

	if cfg.Version {
		return cfg, nil
	}

	if cfg.Hostport == "" {
		return nil, fmt.Errorf("-hostport is required")
	}

	return cfg, nil
}
