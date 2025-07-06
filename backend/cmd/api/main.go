package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"gitlab.praktikum-services.ru/Stasyan/momo-store/cmd/api/app"
	"gitlab.praktikum-services.ru/Stasyan/momo-store/cmd/api/dependencies"
	"gitlab.praktikum-services.ru/Stasyan/momo-store/internal/logger"
)

func main() {
	isHealthCheck := flag.Bool("health-check", false, "Run health check and exit.")
	flag.Parse()

	if *isHealthCheck {
		runHealthCheck()
		return
	}

	logger.Setup()

	if err := run(); err != nil {
		logger.Log.Fatal("unexpected error", zap.Error(err))
		os.Exit(1)
	}
}

func runHealthCheck() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	url := fmt.Sprintf("http://localhost:%s/health", port)

	client := http.Client{
		Timeout: 3 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		logger.Log.Error("health check failed", zap.Error(err))
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Log.Error("health check returned non-200 status", zap.Int("status", resp.StatusCode))
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }
	addr := fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	store, err := dependencies.NewFakeDumplingsStore()
	if err != nil {
		return fmt.Errorf("cannot bootstrap dumplings store: %w", err)
	}

	logger.Log.Debug("creating app instance")
	instance, err := app.NewInstance(store)
	if err != nil {
		return fmt.Errorf("cannot create app instance: %w", err)
	}

	router, err := newRouter(instance)
	if err != nil {
		return fmt.Errorf("cannot create router instance: %w", err)
	}

	srv := &http.Server{
		Handler: router,
	}

	errChan := make(chan error, 1)
	go func() {
		logger.Log.Info("starting HTTP server", zap.String("address", addr))
		if err := srv.Serve(lis); err != nil {
			errChan <- fmt.Errorf("error serving HTTP: %w", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stop:
		logger.Log.Info("shutting down gracefully", zap.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return srv.Shutdown(ctx)
	case err := <-errChan:
		return err
	}
}
