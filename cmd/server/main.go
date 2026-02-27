package main

// @title TaskForge API
// @version 1.0
// @description API для системы управления задачами
// @host localhost:8080
// @BasePath /api/v1

import (
	"TaskForge/cmd/cli"
	"TaskForge/internal/config"
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "TaskForge/docs"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	_ = godotenv.Load()
	cfg := config.LoadConfig()

	// initial configuration
	services := config.NewConfigurations(cfg)
	if services == nil {
		logrus.Fatalf("failed to init services: check config.yaml")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logrus.Info("shutdown signal received")
		cancel()
	}()

	app, err := cli.NewApp(services)
	if err != nil {
		logrus.Fatalf("failed to create app: %v", err)
	}

	serverCtx, serverCancel := context.WithCancel(ctx)
	defer serverCancel()

	go func() {
		if err := app.RunApi(serverCtx); err != nil {
			logrus.Errorf("REST API server stopped: %v", err)
			serverCancel()
		}
	}()

	<-ctx.Done()
	logrus.Info("shutting down REST API server...")

	if services.Redis != nil {
		_ = services.Redis.Close()
	}
	if services.DB != nil {
		_ = services.DB.Close()
	}

	logrus.Info("server exited gracefully")
}
