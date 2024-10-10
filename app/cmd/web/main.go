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
	"log/slog"
	"net/http"
	"os"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

//TODO исправить контексты, кафку, обработать ошибки, исправить все по контракту

var cfg *config.Config
var logger *slog.Logger

func init() {
	cfg = utils.LoadConfig("./config/app.yaml")
	logger = setupLogger(cfg.LoggerConfig.Level)

}

func main() {

	db, err := pkg.NewPsqlClient(context.Background(), cfg)
	if err != nil {
		logger.Error("Error creating db client", "error", err)
	}

	tokenManager := auth.NewTokenService(cfg.JWTConfig.SecretKey, logger)
	emailNotifier := notifier.NewNotifier(logger)

	userRepo := user.NewRepo(db)
	houseRepo := house.NewRepo(db)
	flatRepo := flat.NewRepo(db)
	subscriptionRepo := subscriptions.NewRepo(db)

	//brokers := []string{"localhost:9092"}

	p, err := kafka.NewProducer(cfg.KafkaConfig.Brokers, logger)
	if err != nil {
		logger.Error("Failed to create producer", "error", err)
		return
	}
	defer p.Producer.Close()

	c, err := kafka.NewConsumer(cfg.KafkaConfig.Brokers, subscriptionRepo, emailNotifier, logger)
	if err != nil {
		logger.Error("Failed to create consumer", "error", err)
	}
	defer c.Consumer.Close()

	go func() {
		c.Listen("new-flat")
	}()

	router := mux.NewRouter()

	h := handler.New(tokenManager, userRepo, houseRepo, flatRepo, subscriptionRepo, p, c, logger)

	router.HandleFunc("/register", h.Register).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
	router.HandleFunc("/house/{id}", middleware.Auth(h.CreateHouse, tokenManager, true)).Methods("POST")
	router.HandleFunc("/house/{id}", middleware.Auth(h.GetFlatsByHouseID, tokenManager, true)).Methods("GET")
	router.HandleFunc("/flat/update", middleware.Auth(h.GetModerationFlats, tokenManager, true)).Methods("GET")
	router.HandleFunc("/house/{id}", middleware.Auth(h.GetFlatsByHouseID, tokenManager, false)).Methods("GET")
	router.HandleFunc("/house/{id}/subscribe", middleware.Auth(h.CreateSubscription, tokenManager, false)).Methods("POST")
	router.HandleFunc("/flat/create", middleware.Auth(h.CreateFlat, tokenManager, false)).Methods("POST")

	addr := fmt.Sprintf("%s:%s", cfg.AppCfg.Host, cfg.AppCfg.Port)
	fmt.Println("Server is running on", addr)

	if err = http.ListenAndServe(addr, router); err != nil {
		log.Fatalln("Error launch web server:", err)
	}
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
