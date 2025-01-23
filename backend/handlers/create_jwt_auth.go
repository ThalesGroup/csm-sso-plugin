/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package handlers

import (
	"csm/global"
	"csm/services"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

func CreateJWTAuthMethodHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		CMURL string `json:"cmUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if global.AkeylessToken == "" {
		http.Error(w, "Missing Akeyless t-token", http.StatusUnauthorized)
		return
	}

	jwksEncoded := base64.StdEncoding.EncodeToString([]byte(global.JwksJSON))
	if jwksEncoded == "" {
		http.Error(w, "Missing JWKs JSON", http.StatusNotFound)
		return
	}

	statusCode, err := services.CreateJWTAuthMethod(global.AkeylessToken, requestPayload.CMURL, jwksEncoded)
	if err != nil || statusCode != http.StatusOK {
		http.Error(w, "Failed to create JWT Auth Method: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "JWT Auth created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
