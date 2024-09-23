package main

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/repository/house"
	"app/internal/repository/user"
	"app/pkg"
	"app/pkg/auth"
	"app/pkg/utils"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var cfg *config.Config

func init() {
	cfg = utils.LoadConfig("./config/app.yaml")
}

func main() {
	db, err := pkg.NewPsqlClient(context.Background(), cfg)
	if err != nil {
		log.Fatalln("Error create db client:", err)
	}

	tokenManager := &auth.TokenService{}
	userRepo := user.NewRepo(db)
	houseRepo := house.NewRepo(db)

	router := mux.NewRouter()

	//TODO read
	//h := handler.New(tokenManager, userRepo)

	//err = h.InitRoutes(router, tokenManager, userRepo)
	//if err != nil {
	//	log.Fatalf("Failed to initialize routes: %v", err)
	//}

	h := handler.New(tokenManager, userRepo, houseRepo)
	router.HandleFunc("/register", h.Register).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
	router.HandleFunc("/house/create", h.CreateHouse).Methods("POST")

	port := ":8080"
	fmt.Println("Server is running on", port)
	log.Fatal(http.ListenAndServe(port, router))
}
