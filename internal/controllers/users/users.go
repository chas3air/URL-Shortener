package users

import (
	"URL-Shortener/internal/database/models"
	"log"
	"net/http"
)

type UsersController struct {
	storage *models.DataBase
	client  *http.Client
}

func New(storage *models.DataBase, client *http.Client) *UsersController {
	return &UsersController{
		storage: storage,
		client:  client,
	}
}

func (uc *UsersController) Get(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users (GET) processing")
	log.Println("handler /users (GET) done")
}

func (uc *UsersController) GetById(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users/id (GET) processing")
	log.Println("handler /users/id (GET) done")
}

func (uc *UsersController) GetByLoginAndPassword(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users/login (POST) processing")
	log.Println("handler /users/login (POST) done")
}

func (uc *UsersController) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users (POST) processing")
	log.Println("handler /users (POST) done")
}

func (uc *UsersController) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users/id (PUT) processing")
	log.Println("handler /users/id (PUT) done")
}

func (uc *UsersController) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users/id (DELETE) processing")
	log.Println("handler /users/id (DELETE) done")
}
