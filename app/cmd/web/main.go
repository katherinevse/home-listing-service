package main

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/kafka"
	"app/internal/middleware"
	"app/internal/repository/flat"
	"app/internal/repository/house"
	"app/internal/repository/subscriptions"
	"app/internal/repository/user"
	"app/notifier"
	"app/pkg"
	"app/pkg/auth"
	"app/pkg/utils"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//TODO добавить логгер, middleware исправить контексты, кафку, обработать ошибки, исправить все по контракту

var cfg *config.Config

func init() {
	cfg = utils.LoadConfig("./config/app.yaml")
}

func main() {
	db, err := pkg.NewPsqlClient(context.Background(), cfg)
	if err != nil {
		log.Fatalln("Error create db client:", err)
	}

	tokenManager := auth.NewTokenService(cfg.JWTConfig.SecretKey)
	emailNotifier := &notifier.Notifier{}

	userRepo := user.NewRepo(db)
	houseRepo := house.NewRepo(db)
	flatRepo := flat.NewRepo(db)
	subscriptionRepo := subscriptions.NewRepo(db)

	brokers := []string{"localhost:9092"}

	//Kafka продюсер
	p, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Producer.Close()

	//Kafka консьюмер
	c, err := kafka.NewConsumer(brokers, subscriptionRepo, emailNotifier)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Consumer.Close()

	go func() {
		c.Listen("new-flat")
	}()

	router := mux.NewRouter()

	h := handler.New(tokenManager, userRepo, houseRepo, flatRepo, subscriptionRepo, p, c)

	router.HandleFunc("/register", h.Register).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")

	//проверка на модератора
	router.HandleFunc("/house/{id}", middleware.Auth(h.CreateHouse, tokenManager, true)).Methods("POST")
	router.HandleFunc("/house/{id}", middleware.Auth(h.GetFlatsByHouseID, tokenManager, true)).Methods("GET")
	router.HandleFunc("/flat/update", middleware.Auth(h.GetModerationFlats, tokenManager, true)).Methods("GET")

	//возврат юзера через контекст
	router.HandleFunc("/house/{id}", middleware.Auth(h.GetFlatsByHouseID, tokenManager, false)).Methods("GET")

	//ничего не нужно
	router.HandleFunc("/house/{id}/subscribe", middleware.Auth(h.CreateSubscription, tokenManager, false)).Methods("POST")
	router.HandleFunc("/flat/create", middleware.Auth(h.CreateFlat, tokenManager, false)).Methods("POST")

	addr := fmt.Sprintf("%s:%s", cfg.AppCfg.Host, cfg.AppCfg.Port)
	fmt.Println("Server is running on", addr)

	if err = http.ListenAndServe(addr, router); err != nil {
		log.Fatalln("Error launch web server:", err)
	}
}
