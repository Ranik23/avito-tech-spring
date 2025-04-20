package main

import (
	"context"
	"log"

	"github.com/Ranik23/avito-tech-spring/internal/app"
)

func main() {

	app, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed create App")
	}

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("failed to start")
	}
}