/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package services

import (
	"csm/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// UpdateConfig updates the configuration with the provided details.
func UpdateConfig(token, cmURL, connectionName, ssoAccessID, akeylessURL string) error {
	updateConfigURL := cmURL + "/api/v1/configs/akeyless"

	// Prepare the payload for the PATCH request
	configPayload := map[string]interface{}{
		"gateway_connection_id": connectionName,
		"sso_access_id":         ssoAccessID,
		"akeyless_signup_url":   akeylessURL,
	}

	// Convert the payload to JSON
	body, err := json.Marshal(configPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Set request headers
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	// Attempt to update AkeylessConfig with the PATCH request
	resp, respBody, err := utils.SendRequest("PATCH", updateConfigURL, headers, body)
	if err != nil {
		return err
	}

	// Handle response status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Only print success message if connectionName is not 'None'
		if connectionName != "" {
			fmt.Println("AkeylessConfig updated successfully")
		}
	case http.StatusConflict:
		return fmt.Errorf("conflict occurred while updating config, status code: %d", resp.StatusCode)
	default:
		return fmt.Errorf("failed to update config, status code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
