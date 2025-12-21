package main

import (
	"context"
	"log"

	"github.com/semanggilab/webcore-go/app/config"
	"github.com/semanggilab/webcore-go/app/core"
	"github.com/semanggilab/webcore-go/install"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg := config.Config{}
	// if err := config.LoadConfig(&cfg, "config", "yaml", []string{}); err != nil {
	if err := config.LoadDefaultConfig(&cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize application
	application := core.NewApp(ctx, &cfg, install.APP_LIBRARIES, install.APP_PACKAGES)

	// Start the application
	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
