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

func CheckStatusHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		CMURL string `json:"cmUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use Bearer JWT token
	token := global.BearerToken
	if token == "" {
		http.Error(w, "Missing Bearer token", http.StatusUnauthorized)
		return
	}

	status, err := services.CheckStatus(requestPayload.CMURL, token)
	if err != nil {
		if err.Error() == "status field is missing or invalid" {
			http.Error(w, "CSM tile status is not OK", http.StatusInternalServerError)
		} else {
			http.Error(w, "Failed to check status: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Handle and return the appropriate message based on status
	var message string
	if status == "ready" {
		message = "CSM Tile is enabled & ready to use 😊"
	} else {
		message = "CSM Tile is not enabled yet. Please wait! ☹️"
	}

	response := map[string]interface{}{
		"message": message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
