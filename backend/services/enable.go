/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package services

import (
	"bytes"
	"context"
	"csm/global"
	"csm/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// EnableCSMProduct enables the product in CipherTrust Manager
func EnableCSMProduct(cmURL, token string) error {
	enableProductURL := fmt.Sprintf("%s/api/v1/system/products/%s/enable", cmURL, global.ProductName)

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	// Check the initial service status
	status, err := checkServiceStatus(cmURL, token, global.ServiceName)
	if err != nil {
		return fmt.Errorf("failed to check service status: %w", err)
	}

	switch status {
	case "started":
		fmt.Println("Secrets-Manager service is up and CSM is already enabled.")
		return nil

	case "error":
		fmt.Println("Secrets-Manager service is getting started. Please wait!")

		err := utils.Poll(
			context.Background(),
			func() (interface{}, error) {
				return checkServiceStatus(cmURL, token, global.ServiceName)
			},
			func(result interface{}) bool {
				return result == "started"
			},
			10*time.Minute,
			5*time.Second,
		)
		if err != nil {
			return fmt.Errorf("failed to wait for service to start: %w", err)
		}

		fmt.Println("Secrets-Manager service is now started.")
		return nil

	case "disabled":
		fmt.Println("CSM is disabled in CipherTrust Manager. Enabling it now.")

		resp, _, err := utils.SendRequest("POST", enableProductURL, headers, nil)
		if err != nil {
			return fmt.Errorf("failed to send enable request: %w", err)
		}

		if resp.StatusCode != http.StatusAccepted {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		// Poll for service to be started after enabling the product
		err = utils.Poll(
			context.Background(),
			func() (interface{}, error) {
				return checkServiceStatus(cmURL, token, global.ServiceName)
			},
			func(result interface{}) bool {
				return result == "started"
			},
			10*time.Minute,
			5*time.Second,
		)
		if err != nil {
			return fmt.Errorf("failed to wait for service to start: %w", err)
		}

		fmt.Println("CSM successfully enabled")
		return nil

	default:
		return errors.New("unexpected service status: " + status)
	}
}

// checkServiceStatus checks the secrets-manager service status
func checkServiceStatus(cmURL, token, serviceName string) (string, error) {
	checkServiceStatusURL := fmt.Sprintf("%s/api/v1/system/services/status?service_names=%s", cmURL, serviceName)

	// Set request headers
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	// Send request to get service status
	resp, respBody, err := utils.SendRequest("GET", checkServiceStatusURL, headers, nil)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	// Handle response status codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the JSON response body
	var result map[string]interface{}
	if err := json.NewDecoder(bytes.NewBuffer(respBody)).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	// Assuming result has structure: {"services": [{"status": "ready", "name": "service_name"}]}
	services, ok := result["services"].([]interface{})
	if !ok || len(services) == 0 {
		return "", errors.New("no services found in response")
	}

	// Extract the first service and check its status
	service, ok := services[0].(map[string]interface{})
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	status, ok := service["status"].(string)
	if !ok {
		return "", errors.New("service status not found or invalid")
	}

	return status, nil
}
