# TLS Connection Tester

## Overview and Objectives

This tool is a command-line application written in Go to test TLS connections to a specified host and port. It provides detailed information about the TLS certificate and the negotiated connection parameters.

The main objectives of this tool are:
- To provide a simple way to test TLS connectivity.
- To display detailed certificate information, including Subject, Subject Alternative Name (SAN), Issuer, Root CA, and Expiration Date.
- To show the negotiated TLS Protocol version and Cipher Suite.
- To allow testing of specific TLS versions and cipher suites.
- To support the use of a custom trusted keystore by appending to system certificates.
- To offer an option to bypass chain of trust verification for testing purposes.

## Command-Line Syntax

```sh
tlstester [flags]
```

### Flags

| Flag                 | Description                                                  | Default      |
| -------------------- | ------------------------------------------------------------ | ------------ |
| `-hostport`          | Hostname:port to test (required)                             |              |
| `-timeout`           | Timeout for connection (e.g., `5s`, `10s`)                   | `10s`        |
| `-retries`           | Number of additional connection attempts                      | `3`          |
| `-tls`               | Specific TLS version (`TLS1.0`, `TLS1.1`, `TLS1.2`, `TLS1.3`) |              |
| `-cipher`            | Specific Cipher Suite (e.g., `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`) |              |
| `-keystore`          | Path to PEM-encoded trusted CA certificate file (optional)   |              |
| `-insecure`          | Skip chain of trust verification (use with caution)          | `false`      |
| `-version`           | Display version information                                  | `false`      |

## Usage Examples

### Basic Test

To perform a basic test against a host:
```sh
tlstester -hostport google.com:443
```

### Test a Specific TLS Version

To test the connection forcing TLS 1.2:
```sh
tlstester -hostport google.com:443 -tls TLS1.2
```

### Test a Specific Cipher Suite

To test the connection using a specific cipher suite:
```sh
tlstester -hostport google.com:443 -cipher TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

### Use a Custom Keystore

To use a custom trusted CA certificate (appended to system roots) for verification:
```sh
tlstester -hostport my-internal-host:443 -keystore C:\certs\internal-ca.pem
```

### Skip Certificate Verification

To skip the chain of trust verification (unsafe, for testing only):
```sh
tlstester -hostport self-signed.server:443 -insecure
```

## Troubleshooting

- **Invalid TLS version**: Ensure you use one of `TLS1.0`, `TLS1.1`, `TLS1.2`, or `TLS1.3`.
- **Invalid Cipher Suite**: Use standard Go cipher suite names. Note that TLS 1.3 cipher suites are not configurable via the `-cipher` flag in Go.
- **Connection failed after X retries**: Check network connectivity and ensure the target port is reachable.
- **Failed to parse root certificate**: Ensure the `-keystore` file is a valid PEM-encoded certificate.
