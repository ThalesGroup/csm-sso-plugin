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

func UpdateConfigHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		CMURL          string `json:"cmUrl"`
		ConnectionName string `json:"connection_name"`
		AkeylessURL    string `json:"akeyless_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if global.BearerToken == "" {
		http.Error(w, "Missing Bearer token", http.StatusUnauthorized)
		return
	}

	err := services.UpdateConfig(global.BearerToken, requestPayload.CMURL, requestPayload.ConnectionName, global.SSOAccessID, requestPayload.AkeylessURL)
	if err != nil {
		http.Error(w, "Failed to update config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "AkeylessConfig updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
