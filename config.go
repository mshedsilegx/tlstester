package main

import (
	"flag"
	"fmt"
	"time"
)

type Config struct {
	Hostport           string
	Timeout            time.Duration
	Retries            int
	TLSVersion         string
	CipherSuite        string
	Keystore           string
	InsecureSkipVerify bool
	Version            bool
}

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
