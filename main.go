package main

import (
	"fmt"
	"os"
)

var version string

// main is the entry point for the TLS Connection Tester utility.
// It handles flag parsing, TLS configuration setup, connection establishment,
// and results reporting.
func main() {
	// Parse command-line flags and handle version request
	cfg, err := ParseFlags(version)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if cfg.Version {
		fmt.Printf("TLS Connection Tester - Version: %s\n", version)
		os.Exit(0)
	}

	// Initialize TLS configuration based on user flags
	tlsConfig, err := CreateTLSConfig(cfg)
	if err != nil {
		fmt.Println("Error creating TLS config:", err)
		os.Exit(1)
	}

	// Establish the TLS connection with optional retries
	conn, err := ConnectTLS(cfg, tlsConfig)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing connection: %v\n", err)
		}
	}()

	// Output details about the established connection
	PrintConnectionState(conn)
}
