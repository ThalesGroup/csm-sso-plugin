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

func FetchJWKsHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		CMURL string `json:"cmUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch the JWKs JSON
	err := services.FetchJWKs(requestPayload.CMURL)
	if err != nil {
		// Return error message if fetching JWKs JSON fails
		http.Error(w, "Unable to fetch JWKs JSON from CM: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare success response
	response := map[string]string{
		"message": "Fetched JWKs JSON from CM successfully",
	}

	// Set response headers and encode JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
