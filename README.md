# App Store Connect API SDK for Go

A Go SDK for interacting with the Apple App Store Connect API.

## Features

- JWT authentication with ES256 signing
- Device management (register, list, query by UDID)
- Certificate management (list, create, delete)
- Bundle ID management (register, list, query, delete)
- Profile management (create, list, delete)
- Bundle ID capability management (enable, disable)

## Installation

```bash
go get appstore-connect-api
```

## Configuration

Before using this SDK, you need to obtain API credentials from Apple:

1. Go to [App Store Connect](https://appstoreconnect.apple.com)
2. Navigate to Users and Access > Keys
3. Create a new key and save:
   - Issuer ID
   - Key ID
   - Private Key (.p8 file)

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "appstore-connect-api/pkg/appstore"
)

func main() {
    // Create client configuration
    config := appstore.Config{
        Issuer:    "YOUR_ISSUER_ID",
        KeyID:     "YOUR_KEY_ID",
        Secret:    "./path/to/privatekey.p8", // Path to .p8 file or key content
        APIVersion: "v1",
    }
    
    // Create client
    client, err := appstore.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use Device API
    deviceAPI, _ := client.API("device")
    devices, err := deviceAPI.(*appstore.DeviceAPI).All(nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(devices)
}
```

## API Reference

### Device API

```go
deviceAPI, _ := client.API("device")

// List all devices
devices, err := deviceAPI.(*appstore.DeviceAPI).All(params)

// Register a new device
result, err := deviceAPI.(*appstore.DeviceAPI).Register(name, platform, udid)

// Get device type by UDID
deviceType, _ := deviceAPI.(*appstore.DeviceAPI).GetDeviceType(udid)

// Register device and get type (handles existing devices)
deviceType, _ := deviceAPI.(*appstore.DeviceAPI).RegisterAndGetType(name, platform, udid)

// Get device sort information (available slots)
sortResult, _ := deviceAPI.(*appstore.DeviceAPI).DeviceSort()
```

### Certificates API

```go
certAPI, _ := client.API("certificates")

// List all certificates
certs, err := certAPI.(*appstore.CertificatesAPI).All(params)

// Create a new certificate
newCert, err := certAPI.(*appstore.CertificatesAPI).Create()

// Delete a certificate
result, err := certAPI.(*appstore.CertificatesAPI).Delete(id)
```

### Bundle ID API

```go
bundleIdAPI, _ := client.API("bundleId")

// List all bundle IDs
bundleIds, err := bundleIdAPI.(*appstore.BundleIdAPI).All(params)

// Register a new bundle ID
result, err := bundleIdAPI.(*appstore.BundleIdAPI).Register(name, platform, bundleId)

// Delete a bundle ID
result, err := bundleIdAPI.(*appstore.BundleIdAPI).Delete(bId)

// Query bundle ID capabilities
result, err := bundleIdAPI.(*appstore.BundleIdAPI).Query(bId, params)
```

### Profiles API

```go
profilesAPI, _ := client.API("profiles")

// Query profiles
result, err := profilesAPI.(*appstore.ProfilesAPI).Query(params)

// Create a new profile
result, err := profilesAPI.(*appstore.ProfilesAPI).Create(
    name,
    bId,
    profileType,
    devices,
    certificates,
)

// List devices for a profile
result, err := profilesAPI.(*appstore.ProfilesAPI).ListDevices(pId, params)

// List certificates for a profile
result, err := profilesAPI.(*appstore.ProfilesAPI).ListCertificates(pId, params)

// Delete a profile
result, err := profilesAPI.(*appstore.ProfilesAPI).Delete(pId)
```

### Bundle ID Capability API

```go
capabilityAPI, _ := client.API("bundleIdCapabilities")

// Enable a capability
result, err := capabilityAPI.(*appstore.BundleIdCapabilityAPI).Enable(bId, capability)

// Disable a capability
result, err := capabilityAPI.(*appstore.BundleIdCapabilityAPI).Disable(bcId)
```

## Example

See `examples/main.go` for a complete example demonstrating all API operations.

```bash
cd examples
go run main.go
```

## Project Structure

```
.
├── go.mod
├── pkg/
│   ├── appstore/
│   │   ├── client.go              # Main client
│   │   ├── device.go              # Device API
│   │   ├── certificates.go        # Certificates API
│   │   ├── profiles.go            # Profiles API
│   │   ├── bundleid.go            # Bundle ID API
│   │   └── bundleidcapability.go  # Bundle ID Capability API
│   ├── httpclient/
│   │   └── client.go              # HTTP client
│   └── jwt/
│       └── jwt.go                 # JWT generation
└── examples/
    └── main.go                     # Usage examples
```

## Requirements

- Go 1.21 or higher

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
