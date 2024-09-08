package main

import (
	"go-tender-app/config"
	"go-tender-app/database"
	"go-tender-app/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	database.Connect(cfg.PostgresURL)

	router := mux.NewRouter()

	router.HandleFunc("/api/ping", handlers.Ping).Methods("GET")
	router.HandleFunc("/api/tenders", handlers.GetAllTenders).Methods("GET")
	router.HandleFunc("/api/tenders/new", handlers.CreateTender).Methods("POST")
	router.HandleFunc("/api/tenders/my", handlers.GetUserTenders).Methods("GET")
	router.HandleFunc("/api/tenders/{tenderId}/status", handlers.GetTenderStatus).Methods("GET")
	router.HandleFunc("/api/tenders/{tenderId}/edit", handlers.EditTender).Methods("PATCH")
	router.HandleFunc("/api/tenders/{tenderId}/rollback/{version}", handlers.RollbackTender).Methods("PUT")

	router.HandleFunc("/api/bids/new", handlers.CreateBid).Methods("POST")
	router.HandleFunc("/api/bids/my", handlers.GetUserBids).Methods("GET")
	router.HandleFunc("/api/bids/{bidId}/status", handlers.GetBidStatus).Methods("GET")
	router.HandleFunc("/api/bids/{tenderId}/list", handlers.GetBidsForTender).Methods("GET")
	router.HandleFunc("/api/bids/{bidId}/edit", handlers.EditBid).Methods("PATCH")
	router.HandleFunc("/api/bids/{bidId}/rollback/{version}", handlers.RollbackBidVersion).Methods("PUT")
	router.HandleFunc("/api/bids/{bidId}/submit_decision", handlers.SubmitBidDecision).Methods("PUT")

	http.ListenAndServe(":8080", router)
}
