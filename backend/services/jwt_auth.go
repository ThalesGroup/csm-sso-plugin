/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package services

import (
	"csm/global"
	"csm/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// Creates a new JWT authentication method or handles conflicts if the method already exists.
func CreateJWTAuthMethod(token, cmURL, jwksEncoded string) (int, error) {
	// Define the URL for creating the JWT Auth method
	authURL := cmURL + "/akeyless-api/v2/create-auth-method-oauth2"

	// Prepare the payload for creating the JWT Auth method
	connectionPayload := map[string]interface{}{
		"json":              false,
		"jwks-uri":          "default_jwks_url", // Adjust as necessary
		"jwt-ttl":           60,
		"name":              "JWT-CSM-Plugin",
		"token":             token,
		"unique-identifier": "sub",
		"jwks-json-data":    jwksEncoded,
	}
	body, err := json.Marshal(connectionPayload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal connection payload: %w", err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	// Attempt to create the JWT Auth Method
	resp, respBody, err := utils.SendRequest("POST", authURL, headers, body)
	if err != nil {
		return 0, fmt.Errorf("failed to send request to create JWT Auth Method: %w", err)
	}

	// Handle the response based on status code
	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return resp.StatusCode, fmt.Errorf("failed to decode response body: %w", err)
		}

		if accessID, ok := result["access_id"].(string); ok {
			fmt.Println("JWT Auth successfully created:", accessID)
			// Assuming ssoAccessID is a global or higher-scoped variable that needs to be set
			global.SSOAccessID = accessID
		} else {
			return resp.StatusCode, fmt.Errorf("JWT Access ID not found in response")
		}
	} else if resp.StatusCode == http.StatusConflict {
		// Handle 409 Conflict - delete the existing JWT Auth Method
		deletePayload := map[string]interface{}{
			"json":  false,
			"token": token,
			"name":  "JWT-CSM-Plugin",
		}

		deleteBody, err := json.Marshal(deletePayload)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal delete payload: %w", err)
		}

		deleteURL := cmURL + "/akeyless-api/v2/delete-auth-method"
		resp, respBody, err = utils.SendRequest("POST", deleteURL, headers, deleteBody)
		if err != nil {
			return 0, fmt.Errorf("failed to send request to delete JWT Auth Method: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return resp.StatusCode, fmt.Errorf("failed to delete existing JWT Auth Method, status code: %d, body: %s", resp.StatusCode, string(respBody))
		}

		// Retry creating the JWT Auth Method after deletion
		resp, respBody, err = utils.SendRequest("POST", authURL, headers, body)
		if err != nil {
			return 0, fmt.Errorf("failed to send request to recreate JWT Auth Method: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			if err := json.Unmarshal(respBody, &result); err != nil {
				return resp.StatusCode, fmt.Errorf("failed to decode response body: %w", err)
			}

			if accessID, ok := result["access_id"].(string); ok {
				fmt.Println("JWT Auth successfully deleted and recreated:", accessID)
				global.SSOAccessID = accessID
			} else {
				return resp.StatusCode, fmt.Errorf("JWT Access ID not found in response")
			}
		} else {
			return resp.StatusCode, fmt.Errorf("failed to recreate JWT Auth Method, status code: %d, body: %s", resp.StatusCode, string(respBody))
		}
	} else {
		return resp.StatusCode, fmt.Errorf("unable to create JWT Auth Method, status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return resp.StatusCode, nil
}
