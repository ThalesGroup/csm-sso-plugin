/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package services

import (
	"context"
	"csm/global"
	"csm/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetBearerToken retrieves the bearer and refresh tokens from the authentication server.
func GetBearerToken(cmURL, username, password string) error {
	authURL := cmURL + "/api/v1/auth/tokens/" // Use the correct endpoint
	payload := map[string]string{
		"grant_type": "password",
		"username":   username,
		"password":   password,
		"client_id":  "837c840d-75dd-4b4f-a318-79cb16ca248d",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	// Attempt to get the bearer token
	resp, respBody, err := utils.SendRequest("POST", authURL, headers, body)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Handle response status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Decode the response body
		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}

		// Extract bearer token
		bearerToken, ok := result["jwt"].(string)
		if !ok {
			return fmt.Errorf("missing jwt token in response")
		}

		// Extract refresh token
		refreshTokenID, ok := result["refresh_token_id"].(string)
		if !ok {
			return fmt.Errorf("missing refresh_token in response")
		}

		// Assign tokens to global variables
		global.BearerToken = bearerToken
		global.RefreshTokenID = refreshTokenID

		fmt.Println("Fetched bearer & refresh token ID from CM successfully")
	default:
		return fmt.Errorf("failed to get bearer & refresh token ID, status code: %d", resp.StatusCode)
	}

	return nil
}

// GetAkeylessToken retrieves an Akeyless token using the provided credentials.
func GetAkeylessToken(cmURL, token, akeylessID, akeylessKey string) error {
	authURL := cmURL + "/akeyless-api/v2/auth"

	// Prepare the request payload
	payload := map[string]interface{}{
		"access-type":   "access_key",
		"gcp-audience":  "akeyless.io",
		"json":          false,
		"oci-auth-type": "apikey",
		"access-id":     akeylessID,
		"access-key":    akeylessKey,
	}

	// Marshal the payload to JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Poll to ensure the Secrets-Manager service is started
	err = utils.Poll(
		context.Background(),
		func() (interface{}, error) {
			return checkServiceStatus(cmURL, token, global.ServiceName)
		},
		func(result interface{}) bool {
			return result == "started"
		},
		5*time.Minute, // Polling timeout (total duration)
		5*time.Second, // Frequency of checks
	)
	if err != nil {
		return fmt.Errorf("failed to wait for Secrets-Manager service to start: %w", err)
	}

	fmt.Println("Secrets-Manager service is up. Fetching Akeyless t-token.")

	// Set headers for the request
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	// Send the authentication request
	resp, respBody, err := utils.SendRequest("POST", authURL, headers, body)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Handle response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get Akeyless t-token, status code: %d", resp.StatusCode)
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	// Extract the token from the response
	akeylessToken, ok := response["token"].(string)
	if !ok {
		return fmt.Errorf("akeyless t-token not found in response")
	}

	// Assign the token to a global variable
	global.AkeylessToken = akeylessToken
	fmt.Println("Fetched Akeyless t-token successfully")

	return nil
}

// DeleteToken deletes the provided JWT bearer token.
func DeleteToken(cmURL, refreshTokenID, token string) error {
	deleteTokenURL := cmURL + "/api/v1/auth/tokens/" + refreshTokenID

	// Set headers required for the request
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	// Send the request using the utility function. No body is required for the DELETE request.
	resp, respBody, err := utils.SendRequest("DELETE", deleteTokenURL, headers, nil)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// Handle response status codes
	switch resp.StatusCode {
	case http.StatusNoContent:
		fmt.Println("Token deleted successfully")
	default:
		return fmt.Errorf("failed to delete token, status code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
