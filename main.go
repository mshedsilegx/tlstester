package main

import (
	"fmt"
	"os"
)

var version string

func main() {
	cfg, err := ParseFlags(version)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if cfg.Version {
		fmt.Printf("TLS Connection Tester - Version: %s\n", version)
		os.Exit(0)
	}

	tlsConfig, err := CreateTLSConfig(cfg)
	if err != nil {
		fmt.Println("Error creating TLS config:", err)
		os.Exit(1)
	}

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

	PrintConnectionState(conn)
}
