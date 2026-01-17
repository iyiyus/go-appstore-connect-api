package jwtutil

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtAud = "appstoreconnect-v1"
	jwtAlg = "ES256"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Issuer    string
	KeyID     string
	PrivateKey string
}

// Generator generates JWT tokens for App Store Connect API
type Generator struct {
	config JWTConfig
}

// NewGenerator creates a new JWT generator
func NewGenerator(config JWTConfig) (*Generator, error) {
	if config.Issuer == "" {
		return nil, fmt.Errorf("issuer is required")
	}
	if config.KeyID == "" {
		return nil, fmt.Errorf("key id is required")
	}
	if config.PrivateKey == "" {
		return nil, fmt.Errorf("private key is required")
	}

	return &Generator{config: config}, nil
}

// GenerateToken generates a JWT token
func (g *Generator) GenerateToken() (string, error) {
	// Parse the private key
	privateKey, err := g.parsePrivateKey()
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create token claims
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": g.config.Issuer,
		"iat": now.Add(-60 * time.Second).Unix(), // issued 60 seconds ago
		"exp": now.Add(19 * time.Minute).Unix(),  // expires in 19 minutes
		"aud": jwtAud,
	}

	// Create token with ES256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = g.config.KeyID

	// Sign the token
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// parsePrivateKey parses the private key from string or PEM format
func (g *Generator) parsePrivateKey() (*ecdsa.PrivateKey, error) {
	// Decode PEM block
	block, _ := pem.Decode([]byte(g.config.PrivateKey))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse PKCS8 or PKCS1
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
	}

	// Assert to ECDSA private key
	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not ECDSA")
	}

	return ecdsaKey, nil
}
