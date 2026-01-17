package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Config holds HTTP client configuration
type Config struct {
	BaseURL    string
	APIVersion string
	Token      string
	Headers    map[string]string
}

// Client represents an HTTP client for App Store Connect API
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient creates a new HTTP client
func NewClient(config Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetToken sets the JWT token
func (c *Client) SetToken(token string) {
	c.config.Token = token
}

// SetHeaders sets additional headers
func (c *Client) SetHeaders(headers map[string]string) {
	if c.config.Headers == nil {
		c.config.Headers = make(map[string]string)
	}
	for k, v := range headers {
		c.config.Headers[k] = v
	}
}

// GetHeaders returns all headers including authorization
func (c *Client) GetHeaders() map[string]string {
	headers := make(map[string]string)
	for k, v := range c.config.Headers {
		headers[k] = v
	}
	if c.config.Token != "" {
		headers["Authorization"] = "Bearer " + c.config.Token
	}
	return headers
}

// BuildURL builds the full URL for API requests
func (c *Client) BuildURL(path string) string {
	return fmt.Sprintf("%s/%s%s", c.config.BaseURL, c.config.APIVersion, path)
}

// Get performs a GET request
func (c *Client) Get(path string, params map[string]string) (map[string]interface{}, error) {
	// Build URL
	fullURL := c.BuildURL(path)
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}
		fullURL += "?" + values.Encode()
	}

	// Create request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for k, v := range c.GetHeaders() {
		req.Header.Set(k, v)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return result, nil
}

// PostJSON performs a POST request with JSON body
func (c *Client) PostJSON(path string, body interface{}) (map[string]interface{}, error) {
	// Build URL
	fullURL := c.BuildURL(path)

	// Marshal body
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	headers := c.GetHeaders()
	headers["Content-Type"] = "application/json"
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return result, nil
}

// Delete performs a DELETE request
func (c *Client) Delete(path string, params map[string]string) (map[string]interface{}, error) {
	// Build URL
	fullURL := c.BuildURL(path)
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}
		fullURL += "?" + values.Encode()
	}

	// Create request
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for k, v := range c.GetHeaders() {
		req.Header.Set(k, v)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return result, nil
}
