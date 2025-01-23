/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package services

import (
	"bytes"
	"csm/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// CheckStatus checks the status of the application or service.
func CheckStatus(cmURL, token string) (string, error) {
	statusURL := cmURL + "/api/v1/configs/akeyless/status"

	// Set request headers
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	// Attempt to get the status
	resp, respBody, err := utils.SendRequest("GET", statusURL, headers, nil)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	// Handle response status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Decode the JSON response body
		var result map[string]interface{}
		if err := json.NewDecoder(bytes.NewBuffer(respBody)).Decode(&result); err != nil {
			return "", fmt.Errorf("failed to decode response body: %w", err)
		}

		// Extract the "status" field from the response
		status, ok := result["status"].(string)
		if !ok {
			return "", errors.New("CSM tile status is not OK")
		}

		// Handle the status value
		if status == "ready" {
			fmt.Println("CSM Tile is enabled & ready to use 😊")
		} else {
			fmt.Println("CSM Tile is not enabled yet. Please wait! ☹️")
		}

		return status, nil
	default:
		// Handle unexpected status codes
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
