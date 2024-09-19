package main

import (
	"app/internal/config"
	"app/pkg"
	"app/pkg/utils"
	"context"
	"log"
)

var cfg *config.Config

func init() {
	cfg = utils.LoadConfig("./config/app.yaml")
}

func main() {
	psqlClient, err := pkg.NewPsqlClient(context.Background(), cfg)
	if err != nil {
		log.Fatalln("Error create db client:", err)
	}
}
