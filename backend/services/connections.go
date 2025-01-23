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

// CreateConnection creates a new connection.
func CreateConnection(token, cmURL, connectionName, akeylessID, akeylessKey, akeylessURL string) (int, error) {
	connectionPayload := map[string]interface{}{
		"name":          connectionName,
		"products":      []string{"csm"},
		"access_key_id": akeylessID,
		"access_key":    akeylessKey,
		"akeyless_url":  akeylessURL,
	}

	connectionPayloadBytes, err := json.Marshal(connectionPayload)
	if err != nil {
		return 0, err
	}

	createURL := cmURL + "/api/v1/connectionmgmt/services/akeyless/connections"

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	// Attempt to create the connection
	resp, respBody, err := utils.SendRequest("POST", createURL, headers, connectionPayloadBytes)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == http.StatusConflict {
		// If a 409 Conflict status code is received, reset AkeylessConfig and delete the existing connection
		UpdateConfig(token, cmURL, "", "", "")
		deleteURL := cmURL + "/api/v1/connectionmgmt/services/akeyless/connections/" + connectionName
		resp, respBody, err = utils.SendRequest("DELETE", deleteURL, headers, nil)
		if err != nil {
			return 0, err
		}

		if resp.StatusCode != http.StatusNoContent {
			return resp.StatusCode, fmt.Errorf("failed to delete existing connection, status code: %d, body: %s", resp.StatusCode, string(respBody))
		}

		// Retry creating the connection after deletion
		resp, respBody, err = utils.SendRequest("POST", createURL, headers, connectionPayloadBytes)
		if err != nil {
			return 0, err
		}
	}

	// Check if the status code is 201 OK
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Akeyless connection created successfully")
		return resp.StatusCode, nil
	}

	// Return an error for any other status code
	return resp.StatusCode, fmt.Errorf("failed to create Akeyless connection, status code: %d, body: %s", resp.StatusCode, string(respBody))

}

// TestConnection tests the created connection.
func TestConnection(token, cmURL, connectionName string) (int, error) {
	testURL := cmURL + "/api/v1/connectionmgmt/services/akeyless/connections/" + connectionName + "/test"

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	// Attempt to test the connection
	resp, respBody, err := utils.SendRequest("POST", testURL, headers, nil)
	if err != nil {
		return 0, err
	}

	// Check if the status code is 200 OK
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Akeyless connection tested successfully")
		return resp.StatusCode, nil
	}

	// Return an error for any other status code
	return resp.StatusCode, fmt.Errorf("failed to test connection, status code: %d, body: %s", resp.StatusCode, string(respBody))
}
