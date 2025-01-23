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

func EnableCSMProductHandler(w http.ResponseWriter, r *http.Request) {
	// Struct for decoding the request payload
	var requestPayload struct {
		CMURL string `json:"cmUrl"`
	}

	// Decode the incoming request body
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve the Bearer JWT token from the global config
	token := global.BearerToken
	if token == "" {
		http.Error(w, "Missing Bearer token", http.StatusUnauthorized)
		return
	}

	// Call the service to enable the CSM product
	err := services.EnableCSMProduct(requestPayload.CMURL, token)
	if err != nil {
		http.Error(w, "Failed to enable CSM Product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the success response with the status
	response := map[string]string{
		"message": "CSM Product is enabled",
	}

	// Set response headers and write success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
