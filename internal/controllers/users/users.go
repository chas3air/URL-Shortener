package users

import (
	"URL-Shortener/internal/database/interfaces"
	"URL-Shortener/internal/database/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UsersController struct {
	storage interfaces.Repository[models.User]
	client  *http.Client
}

func New(storage interfaces.Repository[models.User], client *http.Client) *UsersController {
	return &UsersController{
		storage: storage,
		client:  client,
	}
}

func (uc *UsersController) Get(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users (GET) processing")
	users, err := uc.storage.Get()
	if err != nil {
		log.Println("Error fetching users, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(users)
	if err != nil {
		log.Println("Cannot marshal users, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("handler /users (GET) done")
}

func (uc *UsersController) GetById(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users/id/ (GET) processing")
	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("Wrong id:", id_s, ", error:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uc.storage.GetById(id)
	if err != nil {
		log.Println("Error fetching user by id:", id, ", error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(user)
	if err != nil {
		log.Println("Cannot marshal user, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("handler /users/id/ (GET) done")
}

func (uc *UsersController) GetByLoginAndPassword(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users/login (GET) processing")
	login := mux.Vars(r)["login"]
	password := mux.Vars(r)["password"]

	users, err := uc.storage.Get()
	if err != nil {
		log.Println("Error fetching users, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var neededUser models.User
	for _, user := range users {
		if user.Login == login && user.Password == password {
			neededUser = user
			break
		}
	}
	if neededUser.Id == 0 {
		log.Println("User doesn't exist")
		http.Error(w, "User doesn't exist", http.StatusNotFound)
		return
	}

	bs, err := json.Marshal(neededUser)
	if err != nil {
		log.Println("Cannot marshal user, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("handler /users/login (GET) done")
}

func (uc *UsersController) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users (POST) processing")

	var userForInsert models.User
	if err := json.NewDecoder(r.Body).Decode(&userForInsert); err != nil {
		log.Println("Error reading request body, error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := uc.storage.Insert(userForInsert)
	if err != nil {
		log.Println("Error inserting, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Id int `json:"id"`
	}{Id: userID}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response, error:", err)
	}
	log.Println("handler /users (POST) done")
}

func (uc *UsersController) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users/id (PUT) processing")
	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("Wrong id:", id, ", error:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userForUpdate models.User
	if err := json.NewDecoder(r.Body).Decode(&userForUpdate); err != nil {
		log.Println("Error read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = uc.storage.Update(id, userForUpdate)
	if err != nil {
		log.Println("Error update user by id:", id, ", error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	log.Println("handler /users/id (PUT) done")
}

func (uc *UsersController) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /users/id (DELETE) processing")
	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("Wrong id:", id, ", error:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uc.storage.Delete(id)
	if err != nil {
		log.Println("Error deleteing user, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshalling user, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("handler /users/id (DELETE) done")
}
