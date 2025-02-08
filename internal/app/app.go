package app

import (
	"URL-Shortener/internal/config"
	userscontroller "URL-Shortener/internal/controllers/users"
	"URL-Shortener/internal/database/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	cfg     *config.Config
	storage *models.DataBase
	srv     *http.Server
	wg      sync.WaitGroup
}

func New(cfg *config.Config, storage *models.DataBase) *App {
	return &App{
		cfg:     cfg,
		storage: storage,
	}
}

func (a *App) StartServer() error {
	uc := userscontroller.New(a.storage, &http.Client{Timeout: a.cfg.ContextTime})
	_ = uc

	r := mux.NewRouter()
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods(http.MethodGet)
	r.HandleFunc("/users", uc.Get).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}/", uc.GetById).Methods(http.MethodGet)
	r.HandleFunc("/users/login", uc.GetByLoginAndPassword).Methods(http.MethodGet)
	r.HandleFunc("/users", uc.Insert).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", uc.Update).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", uc.Delete).Methods(http.MethodDelete)

	a.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.Port),
		Handler: r,
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("error of starting server", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	return a.Stop()
}

func (a *App) Stop() error {
	log.Println("Stoping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("error while stoping server: %v", err)
	}

	a.wg.Wait()
	log.Println("Server is stoped")
	return nil
}
