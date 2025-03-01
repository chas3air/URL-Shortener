package main

import (
	"URL-Shortener/internal/app"
	"URL-Shortener/internal/config"
	"URL-Shortener/internal/storage/sqlite"

	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log.Print("config was loaded")

	storage := urlrepository.New(cfg.StoragePath)
	log.Println("storage was initialised")

	application := app.New(cfg, storage)

	go func() {
		if err := application.StartServer(); err != nil {
			log.Println("server error:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("stopping application")
	log.Println("application stopped")
}
