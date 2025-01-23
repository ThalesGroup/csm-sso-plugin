/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package handlers

import (
	"csm/services"
	"encoding/json"
	"net/http"
)

func InitializeHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		CMURL    string `json:"cmUrl"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := services.GetBearerToken(requestPayload.CMURL, requestPayload.Username, requestPayload.Password)
	if err != nil {
		http.Error(w, "Failed to get token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the success response
	response := map[string]string{
		"message": "Fetched bearer & refresh token from CM successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
