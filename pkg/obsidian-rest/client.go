package obsidianrest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrNotFound = errors.New("file not found")
)

// Client represents an Obsidian REST API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// ClientOption represents a function that configures a Client
type ClientOption func(*Client)

// WithInsecureSkipVerify configures the client to skip SSL certificate verification
func WithInsecureSkipVerify(skip bool) ClientOption {
	return func(c *Client) {
		if c.httpClient.Transport == nil {
			c.httpClient.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: skip,
				},
			}
		} else if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
			if transport.TLSClientConfig == nil {
				transport.TLSClientConfig = &tls.Config{}
			}
			transport.TLSClientConfig.InsecureSkipVerify = skip
		}
	}
}

// NewClient creates a new Obsidian REST API client
func NewClient(baseURL, apiKey string, opts ...ClientOption) *Client {
	client := &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	// Apply all options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// VaultFile represents a file in the Obsidian vault
type VaultFile struct {
	Content     string                 `json:"content"`
	Frontmatter map[string]interface{} `json:"frontmatter"`
	Path        string                 `json:"path"`
	Stat        FileStat               `json:"stat"`
	Tags        []string               `json:"tags"`
}

// FileStat represents file statistics
type FileStat struct {
	CTime int64 `json:"ctime"`
	MTime int64 `json:"mtime"`
	Size  int64 `json:"size"`
}

// GetVaultFile retrieves a specific file from the vault
func (c *Client) GetVaultFile(path string) (*VaultFile, error) {
	url := fmt.Sprintf("%s/vault/%s", c.baseURL, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Accept", "application/vnd.olrapi.note+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("body", string(body))

	var file VaultFile
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&file); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &file, nil
}

// UpdateVaultFile updates an existing file in the vault
func (c *Client) UpdateVaultFile(path string, content string) error {
	url := fmt.Sprintf("%s/vault/%s", c.baseURL, path)

	req, err := http.NewRequest("PUT", url, bytes.NewBufferString(content))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "text/markdown")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteVaultFile deletes a file from the vault
func (c *Client) DeleteVaultFile(path string) error {
	url := fmt.Sprintf("%s/vault/%s", c.baseURL, path)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
