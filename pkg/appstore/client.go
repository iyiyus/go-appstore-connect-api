package appstore

import (
	"fmt"
	"os"

	"appstore-connect-api/pkg/httpclient"
	"appstore-connect-api/pkg/jwtutil"
)

const (
	baseURI    = "https://api.appstoreconnect.apple.com"
	defaultAPIVersion = "v1"
)

// Config holds the client configuration
type Config struct {
	Issuer    string
	KeyID     string
	Secret    string // Can be a file path or the private key content
	APIVersion string
}

// Client represents the App Store Connect API client
type Client struct {
	config     Config
	httpClient *httpclient.Client
	jwtGenerator *jwtutil.Generator
}

// NewClient creates a new App Store Connect API client
func NewClient(config Config) (*Client, error) {
	// Validate required fields
	if config.Issuer == "" {
		return nil, fmt.Errorf("issuer is required")
	}
	if config.KeyID == "" {
		return nil, fmt.Errorf("key id is required")
	}
	if config.Secret == "" {
		return nil, fmt.Errorf("secret is required")
	}

	// Set default API version
	if config.APIVersion == "" {
		config.APIVersion = defaultAPIVersion
	}

	// Read secret from file if it's a file path
	privateKey := config.Secret
	if _, err := os.Stat(config.Secret); err == nil {
		content, err := os.ReadFile(config.Secret)
		if err != nil {
			return nil, fmt.Errorf("failed to read secret file: %w", err)
		}
		privateKey = string(content)
	}

	// Create JWT generator
	jwtGenerator, err := jwtutil.NewGenerator(jwtutil.JWTConfig{
		Issuer:    config.Issuer,
		KeyID:     config.KeyID,
		PrivateKey: privateKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT generator: %w", err)
	}

	// Create HTTP client
	httpClient := httpclient.NewClient(httpclient.Config{
		BaseURL:    baseURI,
		APIVersion: config.APIVersion,
	})

	return &Client{
		config:      config,
		httpClient:  httpClient,
		jwtGenerator: jwtGenerator,
	}, nil
}

// GetToken generates and returns a JWT token
func (c *Client) GetToken() (string, error) {
	token, err := c.jwtGenerator.GenerateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
}

// EnsureAuth ensures the client has an auth header with JWT token
func (c *Client) EnsureAuth() error {
	if c.httpClient.GetHeaders()["Authorization"] == "" {
		token, err := c.GetToken()
		if err != nil {
			return err
		}
		c.httpClient.SetToken(token)
	}
	return nil
}

// API returns an API client for the specified name
func (c *Client) API(name string) (interface{}, error) {
	switch name {
	case "device":
		return NewDeviceAPI(c), nil
	case "bundleId":
		return NewBundleIdAPI(c), nil
	case "bundleIdCapabilities":
		return NewBundleIdCapabilityAPI(c), nil
	case "profiles":
		return NewProfilesAPI(c), nil
	case "certificates":
		return NewCertificatesAPI(c), nil
	default:
		return nil, fmt.Errorf("undefined API: %s", name)
	}
}

// GetHTTPClient returns the underlying HTTP client
func (c *Client) GetHTTPClient() *httpclient.Client {
	return c.httpClient
}
