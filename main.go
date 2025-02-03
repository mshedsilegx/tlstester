package main

import (
  "crypto/tls"
  "crypto/x509"
  "flag"
  "fmt"
  "net"
  "os"
  "time"
)

func main() {
  hostport := flag.String("hostport", "", "Hostname:port to test (required)")
  timeout := flag.Duration("timeout", 10*time.Second, "Timeout for connection")
  retries := flag.Int("retries", 3, "Number of retries")
  tlsVersion := flag.String("tls", "", "Specific TLS version (e.g., TLS1.2, TLS1.3)")
  cipherSuite := flag.String("cipher", "", "Specific Cipher Suite (e.g., TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256)")
  keystore := flag.String("keystore", "", "Path to trusted keystore (optional)")
  insecureSkipVerify := flag.Bool("insecure", false, "Skip chain of trust verification (use with caution)")

  flag.Parse()

  if *hostport == "" {
    fmt.Println("Error: -hostport is required")
    os.Exit(1)
  }

  var tlsConfig *tls.Config = &tls.Config{
    InsecureSkipVerify: *insecureSkipVerify,
  }

  if !*insecureSkipVerify {  // Only set ServerName if not skipping verification
    hostname, _, err := net.SplitHostPort(*hostport)
    if err != nil {
      fmt.Println("Error parsing hostport:", err)
      os.Exit(1)
    }
    tlsConfig.ServerName = hostname
  }
	
  if *tlsVersion != "" {
    switch *tlsVersion {
      case "TLS1.0": tlsConfig.MinVersion = tls.VersionTLS10; tlsConfig.MaxVersion = tls.VersionTLS10
      case "TLS1.1": tlsConfig.MinVersion = tls.VersionTLS11; tlsConfig.MaxVersion = tls.VersionTLS11
      case "TLS1.2": tlsConfig.MinVersion = tls.VersionTLS12; tlsConfig.MaxVersion = tls.VersionTLS12
      case "TLS1.3": tlsConfig.MinVersion = tls.VersionTLS13; tlsConfig.MaxVersion = tls.VersionTLS13
      default: fmt.Println("Invalid TLS version"); os.Exit(1)
    }
  }

  if *cipherSuite != "" {
    var ciphers []uint16
    for _, suite := range tls.CipherSuites() {
      if suite.Name == *cipherSuite {
        ciphers = append(ciphers, suite.ID)
        break
      }
    }
    if len(ciphers) == 0 {
      fmt.Println("Invalid Cipher Suite"); os.Exit(1)
    }
    tlsConfig.CipherSuites = ciphers
  }

  if *keystore != "" {
    certs, err := loadKeystore(*keystore)
    if err != nil {
      fmt.Println("Error loading keystore:", err)
      os.Exit(1)
    }
    tlsConfig.RootCAs = certs
  }

  for i := 0; i <= *retries; i++ {
    conn, err := net.DialTimeout("tcp", *hostport, *timeout)
    if err != nil {
      if i < *retries {
        fmt.Printf("Connection failed (attempt %d): %v, retrying...\n", i+1, err)
        time.Sleep(1 * time.Second) // Wait before retrying
        continue
      } else {
        fmt.Printf("Connection failed after %d retries: %v\n", *retries+1, err)
        os.Exit(1)
      }
    }
    defer conn.Close()

    tlsConn := tls.Client(conn, tlsConfig)
    defer tlsConn.Close()

    err = tlsConn.Handshake()
    if err != nil {
      fmt.Println("TLS Handshake failed:", err)
      os.Exit(1)
    }

    cert := tlsConn.ConnectionState().PeerCertificates[0]
    fmt.Println("Certificate Information:")
    fmt.Println("  Subject:", cert.Subject)
    fmt.Println("  SAN:", cert.DNSNames)
    fmt.Println("  Issuer:", cert.Issuer)

    if len(tlsConn.ConnectionState().PeerCertificates) > 1 {
      fmt.Println("  Root CA:", tlsConn.ConnectionState().PeerCertificates[len(tlsConn.ConnectionState().PeerCertificates)-1].Subject)
    }

    fmt.Println("  Protocol:", tls.VersionName(tlsConn.ConnectionState().Version))
    fmt.Println("  Cipher Suite:", tls.CipherSuiteName(tlsConn.ConnectionState().CipherSuite))
    return
  }
}

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
