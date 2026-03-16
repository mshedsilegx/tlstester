# TLS Connection Tester - Design Document

## Overview
The **TLS Connection Tester** is a command-line utility written in Go designed to verify TLS connectivity to a specified host and port. It provides detailed information about the established TLS session, including certificate details, negotiated protocol, and cipher suites.

## Specifications

### Core Functionality
- **Target Connection**: Accepts a `host:port` parameter to initiate a connection.
- **Configurable Timeouts**: Users can specify a connection timeout (defaulting to 10 seconds).
- **Retry Mechanism**: Supports multiple connection attempts in case of transient failures.
- **TLS Protocol Selection**: Allows forcing a specific TLS version (TLS 1.0, 1.1, 1.2, or 1.3).
- **Cipher Suite Selection**: Supports specifying a particular cipher suite using the standard Go `crypto/tls` format.
- **Custom Trust Store**: Option to provide a manual path to a trusted keystore (PEM-encoded root certificates).
- **Insecure Mode**: Capability to bypass the chain of trust checks for testing purposes.

### Output Information
Upon a successful handshake, the tool reports:
- **Subject**: The certificate's subject.
- **Subject Alternative Names (SAN)**: List of DNS names associated with the certificate.
- **Issuer**: The certificate issuer.
- **Root CA**: The root certificate in the provided chain.
- **Expiration Date**: Certificate validity end date.
- **Protocol**: The TLS version negotiated (e.g., TLS 1.3).
- **Cipher Suite**: The negotiated cipher suite (e.g., `TLS_AES_128_GCM_SHA256`).

## Design Principles

### 1. Modularity
The application is structured into logical components:
- `main.go`: Entry point handling the high-level execution flow.
- `config.go`: Centralized flag parsing and configuration management using the standard `flag` package.
- `tls_util.go`: Encapsulates TLS configuration logic and connection handling.
- `keystore.go`: Dedicated logic for loading and merging custom root certificates.

### 2. Robustness
- **Error Handling**: Comprehensive error checking during flag parsing, certificate loading, and connection phases.
- **Connection Resilience**: Implements a retry loop for TCP dial operations to handle intermittent network issues.
- **Resource Management**: Uses `defer` to ensure network connections are closed properly.

### 3. Security (for Testing)
- **SNI Support**: Automatically sets the Server Name Indication (SNI) based on the target host unless certificate verification is disabled.
- **System Trust Integration**: Custom keystores are appended to the system's root certificate pool rather than replacing it entirely, maintaining trust in standard CAs while allowing private CA testing.
- **Insecure Option**: Clearly marked as "use with caution" to allow testing of self-signed or expired certificates when necessary.

## Usage Example
```powershell
.\tlstester.exe -hostport example.com:443 -tls TLS1.2 -retries 2
```
