/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package handlers

import (
	"csm/global"
	"csm/services"
	"encoding/json"
	"net/http"
)

func CreateConnectionHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		CMURL          string `json:"cmUrl"`
		ConnectionName string `json:"connection_name"`
		AkeylessID     string `json:"akeyless_id"`
		AkeylessKey    string `json:"akeyless_key"`
		AkeylessURL    string `json:"akeyless_url"`
	}

	// Decode the request payload
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token := global.BearerToken
	if token == "" {
		http.Error(w, "Missing Bearer JWT token", http.StatusUnauthorized)
		return
	}

	// Create the connection using the provided details
	statusCode, err := services.CreateConnection(token, requestPayload.CMURL, requestPayload.ConnectionName, requestPayload.AkeylessID, requestPayload.AkeylessKey, requestPayload.AkeylessURL)
	if err != nil || statusCode != http.StatusCreated {
		http.Error(w, "Failed to create connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Test the connection using the provided details
	statusCode, err = services.TestConnection(token, requestPayload.CMURL, requestPayload.ConnectionName)
	if err != nil || statusCode != http.StatusOK {
		http.Error(w, "Failed to test connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a success response with the test result
	response := map[string]interface{}{
		"message": "Connection created and tested successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
