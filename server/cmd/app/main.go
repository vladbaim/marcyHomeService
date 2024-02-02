package main

import (
	"context"
	"log"
	"marcyHomeService/internal/app"
	"marcyHomeService/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.GetConfig()

	a, err := app.NewApp(cfg)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("running app")
	err = a.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
