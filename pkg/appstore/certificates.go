package appstore

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

// CertificatesAPI handles certificate-related operations
type CertificatesAPI struct {
	client *Client
}

// NewCertificatesAPI creates a new Certificates API client
func NewCertificatesAPI(client *Client) *CertificatesAPI {
	return &CertificatesAPI{client: client}
}

// All retrieves all certificates
func (c *CertificatesAPI) All(params map[string]string) (map[string]interface{}, error) {
	if err := c.client.EnsureAuth(); err != nil {
		return nil, err
	}
	return c.client.GetHTTPClient().Get("/certificates", params)
}

// Delete deletes a certificate by ID
func (c *CertificatesAPI) Delete(id string) (map[string]interface{}, error) {
	if err := c.client.EnsureAuth(); err != nil {
		return nil, err
	}
	return c.client.GetHTTPClient().Delete("/certificates/"+id, nil)
}

// getRandomCSR generates a random Certificate Signing Request
func (c *CertificatesAPI) getRandomCSR() (string, error) {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{"US"},
			Province:           []string{"California"},
			Locality:           []string{"San Francisco"},
			Organization:       []string{"GoAppStore" + randomString(8)},
			OrganizationalUnit: []string{"GoAppStore" + randomString(8)},
			CommonName:         "CommonName" + randomString(8),
		},
		EmailAddresses: []string{"camen" + randomString(8) + "@example.com"},
	}

	// Generate CSR
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create CSR: %w", err)
	}

	// Encode CSR to PEM format
	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	// Remove PEM headers and footers, keep only the base64 content
	csrString := string(csrPEM)
	csrString = pemHeadersToContent(csrString)

	return csrString, nil
}

// Create creates a new certificate
func (c *CertificatesAPI) Create() (map[string]interface{}, error) {
	if err := c.client.EnsureAuth(); err != nil {
		return nil, err
	}

	csrContent, err := c.getRandomCSR()
	if err != nil {
		return nil, fmt.Errorf("failed to generate CSR: %w", err)
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "certificates",
			"attributes": map[string]string{
				"certificateType": "IOS_DISTRIBUTION",
				"csrContent":      csrContent,
			},
		},
	}

	return c.client.GetHTTPClient().PostJSON("/certificates", data)
}

// pemHeadersToContent removes PEM headers and footers from a PEM string
func pemHeadersToContent(pemString string) string {
	content := pemString
	content = trimPrefix(content, "-----BEGIN CERTIFICATE REQUEST-----")
	content = trimPrefix(content, "-----BEGIN CERTIFICATE-----")
	content = trimSuffix(content, "-----END CERTIFICATE REQUEST-----")
	content = trimSuffix(content, "-----END CERTIFICATE-----")
	return trim(content)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[int64(time.Now().UnixNano()+int64(i))%int64(len(charset))]
	}
	return string(b)
}

// Helper functions for string trimming
func trimPrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}

func trimSuffix(s, suffix string) string {
	if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func trim(s string) string {
	return trimSuffix(trimPrefix(s, " "), " ")
}
