package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func gracefulShutdown(cancel context.CancelFunc, done chan struct{}) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-ctx.Done()
	log.Println("graceful shutdown has triggered...")

	time.Sleep(5 * time.Second)
	cancel()

	close(done)
}

func main() {
	done := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())

	go gracefulShutdown(cancel, done)

	<-ctx.Done()
	log.Println("main context canceled, waiting for cleanup...")

	<-done
	log.Println("graceful shutdown has complete.")
}
