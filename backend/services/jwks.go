/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package services

import (
	"csm/global"
	"csm/utils"
	"fmt"
	"net/http"
)

// FetchJWKs fetches the JSON Web Key Set (JWKS) from the provided URL.
func FetchJWKs(cmURL string) error {
	jwksURL := cmURL + "/api/v1/auth/jwks.json" // Correct endpoint as per fetchJWKs

	headers := map[string]string{
		"Accept": "application/json",
	}

	// Use utils.SendRequest to send the request, passing nil for the body since it's a GET request
	resp, respBody, err := utils.SendRequest("GET", jwksURL, headers, nil)
	if err != nil {
		return err
	}

	// Handle response status codes
	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Fetched JWKs JSON successfully")
	default:
		return fmt.Errorf("failed to fetch JWKs JSON, status code: %d, response: %s", resp.StatusCode, string(respBody))
	}
	global.JwksJSON = string(respBody)

	return nil
}
