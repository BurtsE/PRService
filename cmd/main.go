package main

import (
	"PRService/internal/app"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	a := app.App{}
	logger := a.Logger()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.Start(context.Background()); err != nil {
			logger.Errorf("Failed to start server: %v", err)
			sigchan <- syscall.SIGTERM
		}
	}()

	<-sigchan

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := a.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("Graceful shutdown failed: %v", err)
	}
}
