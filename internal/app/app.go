package app

import (
	"URL-Shortener/internal/config"
	urlcontroller "URL-Shortener/internal/controllers/url"
	"URL-Shortener/internal/storage"

	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type App struct {
	cfg     *config.Config
	storage storage.URLRepository
	srv     *http.Server
	wg      sync.WaitGroup
}

func New(cfg *config.Config, storage storage.URLRepository) *App {
	return &App{
		cfg:     cfg,
		storage: storage,
	}
}

func (a *App) StartServer() error {
	r := mux.NewRouter()
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods(http.MethodGet)

	rc := urlcontroller.New(a.storage, &http.Client{Timeout: a.cfg.ContextTime})
	r.HandleFunc("/url", rc.Get).Methods(http.MethodGet)
	r.HandleFunc("/url/{alias}", rc.GetByAlias).Methods(http.MethodGet)
	r.HandleFunc("/url", rc.Insert).Methods(http.MethodPost)
	r.HandleFunc("/url/{alias}", rc.Delete).Methods(http.MethodDelete)

	a.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.Port),
		Handler: r,
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.srv.ListenAndServe(); err != nil {
			log.Println("error of starting server", "error", err)
		}
	}()

	return nil
}
