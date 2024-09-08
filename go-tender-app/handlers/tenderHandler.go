package handlers

import (
	"encoding/json"
	"go-tender-app/models"
	"go-tender-app/service/tenderservice"

	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllTenders(w http.ResponseWriter, r *http.Request) {
	rows, err := tenderservice.GetAllTenders()

	if err != nil {
		log.Printf("Error querying tenders: %v", err)
		http.Error(w, "Unable to retrieve tenders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tenders, err := tenderservice.Validate(&rows)
	if err != nil {
		log.Printf("Error scanning tender: %v", err)
		http.Error(w, "Error processing tenders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenders)
}

func CreateTender(w http.ResponseWriter, r *http.Request) {
	var tender models.Tender
	err := json.NewDecoder(r.Body).Decode(&tender)
	if err != nil {
		log.Printf("Error fetching user tenders: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = tenderservice.Create(&tender)

	if err != nil {
		log.Printf("Error creating tender: %v", err)
		http.Error(w, "Unable to create tender", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tender)
}

func GetUserTenders(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("username")

	if userName == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	rows, err := tenderservice.GetByUsername(&userName)

	if err != nil {
		log.Printf("Error fetching user tenders: %v", err)
		http.Error(w, "Unable to fetch user tenders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tenders, err := tenderservice.Validate(&rows)
	if err != nil {
		log.Printf("Error scanning user tender: %v", err)
		http.Error(w, "Error processing user tenders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenders)
}

func GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderIdStr := vars["tenderId"]
	tenderID, err := strconv.Atoi(tenderIdStr)

	if err != nil {
		http.Error(w, "Invalid tenderId", http.StatusBadRequest)
		return
	}
	status, err := tenderservice.GetStatus(tenderID)
	if err != nil {
		log.Printf("Error querying tender status: %v", err)
		http.Error(w, "Unable to retrieve tender status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]models.TenderStatus{"status": *status})
}

func EditTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID, err := strconv.Atoi(vars["tenderId"])
	if err != nil {
		http.Error(w, "Invalid tenderId", http.StatusBadRequest)
		return
	}
	var tender models.Tender
	err = json.NewDecoder(r.Body).Decode(&tender)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err, _ = tenderservice.UpdateTender(&tender, tenderID)

	if err != nil {
		log.Printf("Error updating tender: %v", err)
		http.Error(w, "Unable to update tender", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RollbackTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderIdStr := vars["tenderId"]
	versionStr := vars["version"]

	tenderId, err := strconv.Atoi(tenderIdStr)
	if err != nil {
		http.Error(w, "Invalid tenderId", http.StatusBadRequest)
		return
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	newTender, err := tenderservice.RollbackTender(tenderId, version)

	if err != nil {
		log.Printf("Error can't rollback tender: %v", err)
		http.Error(w, "Error can't rollback tender", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newTender)
}

func PublishTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderIDStr := vars["tenderId"]

	tenderID, err := strconv.Atoi(tenderIDStr)
	if err != nil {
		http.Error(w, "Invalid tender ID", http.StatusBadRequest)
		return
	}

	_, err = tenderservice.Publish(tenderID)
	if err != nil {
		http.Error(w, "Error publishing tender", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CloseTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderIDStr := vars["tenderId"]

	tenderID, err := strconv.Atoi(tenderIDStr)
	if err != nil {
		http.Error(w, "Invalid tender ID", http.StatusBadRequest)
		return
	}

	_, err = tenderservice.Close(tenderID)
	if err != nil {
		http.Error(w, "Error closing tender", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
