/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

// SendRequest sends an HTTP request with the given method, URL, headers, and body.
func SendRequest(method, url string, headers map[string]string, body []byte) (*http.Response, []byte, error) {
	// Ensure URL is not empty
	if url == "" {
		return nil, nil, fmt.Errorf("URL cannot be empty")
	}

	// Create a custom HTTP client with TLS configuration that skips certificate verification
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // This will skip TLS certificate verification
			},
		},
	}

	// Create a new HTTP request with body if it's not nil
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, nil, err
	}

	// Set the headers for the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request using the custom HTTP client
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	// Ensure that we defer closing the body only if resp is not nil
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, respBody, nil
}
