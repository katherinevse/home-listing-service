package main

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/kafka"
	"app/internal/repository/flat"
	"app/internal/repository/house"
	"app/internal/repository/subscriptions"
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

//TODO добавить логгер, исправить контексты, кафку, обработать ошибки, исправить все по контракту

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
	subcriptionRepo := subscriptions.NewRepo(db)

	brokers := []string{"localhost:9092"}

	//Kafka продюсер
	p, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Producer.Close()

	//Kafka консьюмер
	c, err := kafka.NewConsumer(brokers, subcriptionRepo)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Consumer.Close()

	go func() {
		c.Listen("new-flat")
	}()

	router := mux.NewRouter()

	h := handler.New(tokenManager, userRepo, houseRepo, flatRepo, subcriptionRepo, p, c)
	router.HandleFunc("/register", h.Register).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
	router.HandleFunc("/house/create", h.CreateHouse).Methods("POST")
	router.HandleFunc("/flat/create", h.CreateFlat).Methods("POST")
	router.HandleFunc("/flat/update", h.GetModerationFlats).Methods("GET")
	router.HandleFunc("/house/{id}", h.GetFlatsByHouseID).Methods("GET")
	router.HandleFunc("/house/{id}/subscribe", h.CreateSubscription).Methods("POST")

	port := ":8080"
	fmt.Println("Server is running on", port)
	log.Fatal(http.ListenAndServe(port, router))
}
