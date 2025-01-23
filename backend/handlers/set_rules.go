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

func SetRulesForAccessRoleHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		CMURL string `json:"cmUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if global.AkeylessToken == "" {
		http.Error(w, "Missing Akeyless token", http.StatusUnauthorized)
		return
	}

	err := services.SetRulesForAccessRole(global.AkeylessToken, requestPayload.CMURL)
	if err != nil {
		http.Error(w, "Failed to set rules for access role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Rules for Access role set successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
