package records

import (
	"URL-Shortener/internal/database/interfaces"
	"URL-Shortener/internal/database/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type RecordsController struct {
	storage interfaces.Repository[models.DbRecord]
	client  *http.Client
}

func New(storage interfaces.Repository[models.DbRecord], client *http.Client) *RecordsController {
	return &RecordsController{
		storage: storage,
		client:  client,
	}
}

func (rc *RecordsController) Get(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /records (GET) processing")
	records, err := rc.storage.Get()
	if err != nil {
		log.Println("Error fetching records, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(records)
	if err != nil {
		log.Println("Cannot marshal records, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("handler /records (GET) done")
}

func (rc *RecordsController) GetById(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /records/{id}/ (GET) processing")
	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("Wrong id:", id_s, ", error:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := rc.storage.GetById(id)
	if err != nil {
		log.Println("Error fetching record by id:", id, ", error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(record)
	if err != nil {
		log.Println("Cannot marshal record, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("handler /records/{id}/ (GET) done")
}

func (rc *RecordsController) GetByAlias(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /records/{alias} (GET) processing")
	alias := mux.Vars(r)["id"]

	records, err := rc.storage.Get()
	if err != nil {
		log.Println("Error fetching records, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var neededRecord models.DbRecord
	for _, record := range records {
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
	log.Println("handler /records/{alias} (GET) done")
}

func (rc *RecordsController) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /records (POST) processing")

	var recordForInsert models.DbRecord
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
	log.Println("handler /records (POST) done")
}

func (rc *RecordsController) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /records/{id} (PUT) processing")
	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("Wrong id:", id, ", error:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var recordForUpdate models.DbRecord
	if err := json.NewDecoder(r.Body).Decode(&recordForUpdate); err != nil {
		log.Println("Error read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rc.storage.Update(id, recordForUpdate)
	if err != nil {
		log.Println("Error update record by id", id, ", error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	log.Println("handler /records/{id} (PUT) done")
}

func (rc *RecordsController) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /records/{id} (DELETE) processing")
	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("Wrong id:", id, ", error:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := rc.storage.Delete(id)
	if err != nil {
		log.Println("Error deleting record, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(record)
	if err != nil {
		log.Println("Error marshalling record, error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("handler /records/{id} (DELETE) done")
}
