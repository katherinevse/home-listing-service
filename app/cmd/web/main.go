package main

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/kafka"
	"app/internal/repository/flat"
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

//TODO добавить логгер, исправить контексты, кафку, обработать ошибки

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
	flatRepo := flat.NewRepo(db)

	//Kafka продюсер
	p, err := kafka.NewProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Producer.Close()

	//Kafka консьюмер

	router := mux.NewRouter()

	h := handler.New(tokenManager, userRepo, houseRepo, flatRepo, p)
	router.HandleFunc("/register", h.Register).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
	router.HandleFunc("/house/create", h.CreateHouse).Methods("POST")
	router.HandleFunc("/flat/create", h.CreateFlat).Methods("POST")
	router.HandleFunc("/flat/update", h.GetModerationFlats).Methods("GET")
	router.HandleFunc("/house/{id}", h.GetFlatsByHouseID).Methods("GET")

	port := ":8080"
	fmt.Println("Server is running on", port)
	log.Fatal(http.ListenAndServe(port, router))
}
