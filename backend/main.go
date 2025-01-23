/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package main

import (
	"csm/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/initialize", handlers.InitializeHandler)
	http.HandleFunc("/csm-service", handlers.EnableCSMProductHandler)
	http.HandleFunc("/create-connection", handlers.CreateConnectionHandler)
	http.HandleFunc("/fetch-jwks", handlers.FetchJWKsHandler)
	http.HandleFunc("/auth-akeyless", handlers.AuthAkeylessHandler)
	http.HandleFunc("/create-jwt-auth", handlers.CreateJWTAuthMethodHandler)
	http.HandleFunc("/create-access-role", handlers.CreateAccessRoleHandler)
	http.HandleFunc("/set-role-rule", handlers.SetRulesForAccessRoleHandler)
	http.HandleFunc("/associate-role", handlers.AssociateAccessRoleHandler)
	http.HandleFunc("/update-config", handlers.UpdateConfigHandler)
	http.HandleFunc("/check-status", handlers.CheckStatusHandler)
	http.HandleFunc("/delete-token", handlers.DeleteTokenHandler)

	fmt.Println("Server starting on port 52920...")
	if err := http.ListenAndServe(":52920", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
