package handler

import (
	"app/pkg/auth"
	"fmt"
	"github.com/gorilla/mux"
)

func InitRoutes(router *mux.Router, tokenManager auth.TokenManager, userRepo UserRepository) error {
	handler := New(tokenManager, userRepo)
	err := router.HandleFunc("/api/register", handler.Register).Methods("POST")
	if err != nil {
		return fmt.Errorf("failed to register route: %w", err)
	}

	//fmt.Printf("Registered routes:\n%s\n", router.Methods("POST").PathPrefix("/"))

	return nil
}
