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

// Set rules for access role for multiple rule types
func SetRulesForAccessRole(token, cmURL string) error {
	// Define the URL for setting rules for access role
	setRuleURL := cmURL + "/akeyless-api/v2/set-role-rule"

	capabilities := []string{"read", "create", "update", "delete", "list"}

	// List of all possible rule types
	ruleTypes := []string{
		"item-rule", "target-rule", "role-rule", "auth-method-rule",
	}

	// Loop through each rule type, construct payload, and send the request
	for _, ruleType := range ruleTypes {
		// Construct the payload for the current rule type
		setRulePayload := map[string]interface{}{
			"json":       false,
			"capability": capabilities,
			"role-name":  "sso-auth",
			"path":       "/*",
			"token":      token,
			"rule-type":  ruleType,
		}

		// Convert the payload to JSON
		body, err := json.Marshal(setRulePayload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload for ruleType %s: %w", ruleType, err)
		}

		// Prepare the headers for the request
		headers := map[string]string{
			"Authorization": "Bearer " + token,
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		}

		// Send the request using the existing utility function
		resp, _, err := utils.SendRequest("POST", setRuleURL, headers, body)
		if err != nil {
			return fmt.Errorf("failed to send request for ruleType %s: %w", ruleType, err)
		}

		// Handle response status codes
		switch resp.StatusCode {
		case http.StatusOK:
			fmt.Printf("Rules for ruleType %s set successfully\n", ruleType)
		default:
			return fmt.Errorf("failed to set rules for ruleType %s, status code: %d", ruleType, resp.StatusCode)
		}
	}

	return nil
}
