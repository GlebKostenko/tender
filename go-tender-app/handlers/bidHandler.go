package handlers

import (
	"encoding/json"
	"go-tender-app/models"
	"go-tender-app/service/bidservice"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateBid(w http.ResponseWriter, r *http.Request) {
	var bid models.Bid
	err := json.NewDecoder(r.Body).Decode(&bid)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = bidservice.Create(&bid)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating bid", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bid)
}

func GetUserBids(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	rows, err := bidservice.GetByUserName(&username)
	if err != nil {
		http.Error(w, "Error fetching bids", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		if err := rows.Scan(&bid.ID, &bid.TenderID, &bid.Description, &bid.Status, &bid.Version, &bid.CreatedAt, &bid.UpdatedAt, &bid.BidID, &bid.UserName, &bid.Name); err != nil {
			log.Println(err)
			http.Error(w, "Error scanning bid", http.StatusInternalServerError)
			return
		}
		bids = append(bids, bid)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bids)
}

func GetBidsForTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderIDStr := vars["tenderId"]

	tenderId, err := strconv.Atoi(tenderIDStr)
	if err != nil {
		http.Error(w, "Invalid tender ID", http.StatusBadRequest)
		return
	}

	rows, err := bidservice.GetBidsByTenderId(tenderId)
	if err != nil {
		http.Error(w, "Error fetching bids", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		if err := rows.Scan(&bid.ID, &bid.TenderID, &bid.Description, &bid.Status, &bid.Version, &bid.CreatedAt, &bid.UpdatedAt, &bid.BidID, &bid.UserName, &bid.Name); err != nil {
			http.Error(w, "Error scanning bid", http.StatusInternalServerError)
			return
		}
		bids = append(bids, bid)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bids)
}

func EditBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidIDStr := vars["bidId"]

	bidId, err := strconv.Atoi(bidIDStr)
	if err != nil {
		http.Error(w, "Invalid bid ID", http.StatusBadRequest)
		return
	}

	var update models.Bid
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err, _ = bidservice.UpdateBid(&update, bidId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating bid", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RollbackBidVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidIDStr := vars["bidId"]
	versionStr := vars["version"]

	bidId, err := strconv.Atoi(bidIDStr)
	if err != nil {
		http.Error(w, "Invalid bid ID", http.StatusBadRequest)
		return
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	newBid, err := bidservice.RollbackBid(bidId, version)
	if err != nil {
		http.Error(w, "Error rolling back bid", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newBid)
}

func GetBidStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidIdStr := vars["bidId"]
	bidID, err := strconv.Atoi(bidIdStr)

	if err != nil {
		http.Error(w, "Invalid tenderId", http.StatusBadRequest)
		return
	}
	status, err := bidservice.GetStatus(bidID)
	if err != nil {
		log.Printf("Error querying tender status: %v", err)
		http.Error(w, "Unable to retrieve tender status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]models.BidStatus{"status": *status})
}

func SubmitBidDecision(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidIDStr := vars["bidId"]

	bidId, err := strconv.Atoi(bidIDStr)
	if err != nil {
		http.Error(w, "Invalid bid ID", http.StatusBadRequest)
		return
	}

	var decision struct {
		Decision string `json:"decision"`
	}
	if err := json.NewDecoder(r.Body).Decode(&decision); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = bidservice.MakeDecision(decision.Decision, bidId)
	if err != nil {
		http.Error(w, "Error submitting decision", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
