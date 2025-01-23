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

func AuthAkeylessHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		CMURL       string `json:"cmUrl"`
		AkeylessID  string `json:"akeyless_id"`
		AkeylessKey string `json:"akeyless_key"`
	}

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

	// Fetch the Akeyless token
	err := services.GetAkeylessToken(requestPayload.CMURL, token, requestPayload.AkeylessID, requestPayload.AkeylessKey)
	if err != nil {
		// Return error message if token fetch fails
		http.Error(w, "Unable to fetch t-token from Akeyless: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare success response with the token
	response := map[string]string{
		"message": "Akeyless t-token fetched successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
