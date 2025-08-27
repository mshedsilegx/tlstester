# TLS Connection Tester

## Overview and Objectives

This tool is a command-line application written in Go to test TLS connections to a specified host and port. It provides detailed information about the TLS certificate and the negotiated connection parameters.

The main objectives of this tool are:
- To provide a simple way to test TLS connectivity.
- To display detailed certificate information, including Subject, Subject Alternative Name (SAN), Issuer, and Root CA.
- To allow testing of specific TLS versions and cipher suites.
- To support the use of a custom trusted keystore.
- To offer an option to bypass chain of trust verification for testing purposes.

## Command-Line Syntax

```
go run . [flags]
```

### Flags

| Flag                 | Description                                                  | Default      |
| -------------------- | ------------------------------------------------------------ | ------------ |
| `-hostport`          | Hostname:port to test (required)                             |              |
| `-timeout`           | Timeout for connection                                       | `10s`        |
| `-retries`           | Number of retries                                            | `3`          |
| `-tls`               | Specific TLS version (e.g., `TLS1.2`, `TLS1.3`)              |              |
| `-cipher`            | Specific Cipher Suite (e.g., `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`) |              |
| `-keystore`          | Path to trusted keystore (optional)                          |              |
| `-insecure`          | Skip chain of trust verification (use with caution)          | `false`      |
| `-version`           | Display version information                                  | `false`      |

## Usage Examples

### Basic Test

To perform a basic test against a host:
```sh
go run . -hostport google.com:443
```

### Test a Specific TLS Version

To test the connection using TLS 1.2:
```sh
go run . -hostport google.com:443 -tls TLS1.2
```

### Test a Specific Cipher Suite

To test the connection using a specific cipher suite:
```sh
go run . -hostport google.com:443 -cipher TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

### Use a Custom Keystore

To use a custom trusted keystore for verification:
```sh
go run . -hostport my-internal-host:443 -keystore /path/to/my/ca.pem
```

### Skip Certificate Verification

To skip the chain of trust verification (unsafe, for testing only):
```sh
go run . -hostport self-signed.server:443 -insecure
```
