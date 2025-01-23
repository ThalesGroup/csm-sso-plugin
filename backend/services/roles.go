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

// Create an access role with a resource.
func CreateAccessRole(token, cmURL string) error {
	// Define the URL for creating access role
	createRoleURL := cmURL + "/akeyless-api/v2/create-role"

	// Prepare the payload for the POST request
	createRolePayload := map[string]interface{}{
		"json":                false,
		"analytics-access":    "own",
		"audit-access":        "own",
		"gw-analytics-access": "own",
		"event-center-access": "own",
		"description":         "Access Role created by CSM Plugin",
		"name":                "sso-auth",
		"token":               token,
	}

	// Convert the payload to JSON
	body, err := json.Marshal(createRolePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Prepare the headers for the request
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	// Send the request using the existing utility function
	resp, _, err := utils.SendRequest("POST", createRoleURL, headers, body)
	if err != nil {
		return fmt.Errorf("failed to send request to create access role: %w", err)
	}

	// Handle response status codes
	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Access Role created successfully")
	case http.StatusConflict:
		// Define the URL for deleting the existing access role
		deleteRoleURL := cmURL + "/akeyless-api/v2/delete-role"

		// Prepare the payload for deleting the existing access role
		deleteRolePayload := map[string]interface{}{
			"json":  false,
			"name":  "sso-auth",
			"token": token,
		}

		// Convert the delete payload to JSON
		deleteBody, err := json.Marshal(deleteRolePayload)
		if err != nil {
			return fmt.Errorf("failed to marshal delete payload: %w", err)
		}

		// Send the delete request
		deleteResp, _, err := utils.SendRequest("POST", deleteRoleURL, headers, deleteBody)
		if err != nil {
			return fmt.Errorf("failed to send request to delete access role: %w", err)
		}

		// Handle the delete response
		if deleteResp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to delete existing access role, status code: %d", deleteResp.StatusCode)
		}

		// Recreate access role
		resp, _, err := utils.SendRequest("POST", createRoleURL, headers, body)
		if err != nil {
			return fmt.Errorf("failed to send request to create access role: %w", err)
		}

		// Check the response for the recreated access role
		if resp.StatusCode == http.StatusOK {
			fmt.Println("Access Role successfully deleted and recreated")
		} else {
			return fmt.Errorf("failed to recreate access role, status code: %d", resp.StatusCode)
		}
	default:
		return fmt.Errorf("failed to create access role, status code: %d", resp.StatusCode)
	}

	return nil
}

// Associates an access role with a resource.
func AssociateAccessRole(token, cmURL string) error {
	// Define the URL for associating access role with Auth Method
	associateRoleURL := cmURL + "/akeyless-api/v2/assoc-role-am"

	// Prepare the payload for the POST request
	associateRolePayload := map[string]interface{}{
		"json":      false,
		"role-name": "sso-auth",
		"am-name":   "JWT-CSM-Plugin",
		"token":     token,
	}

	// Convert the payload to JSON
	body, err := json.Marshal(associateRolePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Prepare the headers for the request
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	// Send the request using the existing utility function
	resp, _, err := utils.SendRequest("POST", associateRoleURL, headers, body)
	if err != nil {
		return fmt.Errorf("failed to send request to associate access role: %w", err)
	}

	// Handle response status codes
	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Access Role successfully associated to Auth Method")
	case http.StatusConflict:
		fmt.Println("Role already associated with the Auth Method")
	default:
		return fmt.Errorf("failed to associate access role to Auth Method, status code: %d", resp.StatusCode)
	}

	return nil
}
