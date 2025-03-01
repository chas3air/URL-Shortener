package url

import (
	"URL-Shortener/internal/models"
	"URL-Shortener/internal/storage"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type URLController struct {
	storage storage.URLRepository
	client  *http.Client
}

func New(storage storage.URLRepository, client *http.Client) *URLController {
	return &URLController{
		storage: storage,
		client:  client,
	}
}

func (rc *URLController) Get(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /url (GET) processing")
	url, err := rc.storage.Get()
	if err != nil {
		log.Println("Error fetching url, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(url)
	if err != nil {
		log.Println("Cannot marshal url, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("handler /url (GET) done")
}

func (rc *URLController) GetByAlias(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /url/{alias} (GET) processing")
	alias := mux.Vars(r)["alias"]

	url, err := rc.storage.Get()
	if err != nil {
		log.Println("Error fetching url, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var neededRecord models.URL
	for _, record := range url {
		if record.Alias == alias {
			neededRecord = record
			break
		}
	}

	if neededRecord.Id == 0 {
		log.Println("Record doesn't exist")
		http.Error(w, "Record doesn't exist", http.StatusNotFound)
		return
	}

	bs, err := json.Marshal(neededRecord)
	if err != nil {
		log.Println("Cannot marshal record, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("handler /url/{alias} (GET) done")
}

func (rc *URLController) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /url (POST) processing")

	var recordForInsert models.URL
	if err := json.NewDecoder(r.Body).Decode(&recordForInsert); err != nil {
		log.Println("Error reading request body, error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recordID, err := rc.storage.Insert(recordForInsert)
	if err != nil {
		log.Println("Error inserting, error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Id int `json:"id"`
	}{Id: recordID}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response, error:", err)
	}
	log.Println("handler /url (POST) done")
}

func (rc *URLController) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /url/{alias} (DELETE) processing")
	alias := mux.Vars(r)["alias"]

	err := rc.storage.Delete(alias)
	if err != nil {
		log.Println("Error deleting record, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	log.Println("handler /url/{id} (DELETE) done")
}
