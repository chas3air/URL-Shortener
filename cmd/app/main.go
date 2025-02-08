package main

import (
	"URL-Shortener/internal/app"
	"URL-Shortener/internal/config"
	"URL-Shortener/internal/database"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log.Print("config was loaded")

	storage := database.MustGetInstanseOfDatabase(cfg.StoragePath)
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
	if err := application.Stop(); err != nil {
		log.Println("error stopping application:", err)
	}
	log.Println("application stopped")
}
