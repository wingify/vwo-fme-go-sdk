/**
 * Copyright 2024 Wingify Software Pvt. Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// Client represents an HTTP client with a base URL and a standard HTTP client.
type Client struct {
	BaseURL string
	client  *http.Client
}

var (
	instance *Client
	once     sync.Once
)

// NewClient creates a new Client instance with the provided base URL.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		client:  &http.Client{},
	}
}

// InitializeClient initializes the singleton instance of Client with the provided base URL.
// It ensures that only one instance is created using sync.Once.
func InitializeClient(baseURL string) {
	once.Do(func() {
		instance = NewClient(baseURL)
	})
}

// GetClient returns the singleton instance of Client.
// It panics if the client has not been initialized.
func GetClient() *Client {
	if instance == nil {
		panic("HTTP client not initialized. Call InitializeClient first.")
	}
	return instance
}

// DoRequest performs an HTTP request with the given method, endpoint, headers, and body.
// It returns the response body or an error if the request fails.
func (c *Client) DoRequest(method, endpoint string, headers map[string]string, body []byte) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating new HTTP request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("received error response from server: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return respBody, nil
}

// PostRequest makes a POST request to the specified endpoint with the given body.
// It sets the Content-Type header to application/json and returns the response body or an error if the request fails.
func (c *Client) PostRequest(endpoint string, body []byte) ([]byte, error) {
	// Headers for the request
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	return c.DoRequest(http.MethodPost, endpoint, headers, body)
}
