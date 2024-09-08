package bidservice

import (
	"context"
	"go-tender-app/database"
	"go-tender-app/models"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

func Create(bid *models.Bid) error {
	var maxBidID int
	err := database.DataSource.QueryRow(context.Background(), "SELECT COALESCE(MAX(bid_id), 0) FROM bids").Scan(&maxBidID)
	if err != nil {
		return err
	}

	bidID := maxBidID + 1

	return database.DataSource.QueryRow(context.Background(),
		`INSERT INTO bids (bid_id, tender_id, name, description, username, status, version)  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		bidID, bid.TenderID, bid.Name, bid.Description, bid.UserName, models.BidCreated, 1).Scan(&bid.ID, &bid.CreatedAt, &bid.UpdatedAt)
}

func GetStatus(bidId int) (*models.BidStatus, error) {
	var status models.BidStatus
	err := database.DataSource.QueryRow(context.Background(), `SELECT status FROM bids WHERE bid_id=$1 AND version=(SELECT MAX(version) FROM tenders WHERE bid_id=$1)`, bidId).Scan(&status)
	return &status, err
}

func GetByUserName(userName *string) (pgx.Rows, error) {
	return database.DataSource.Query(context.Background(), `SELECT * FROM bids WHERE username=$1`, userName)
}

func UpdateBid(update *models.Bid, bidId int) (pgconn.CommandTag, error, int) {
	currentBid, err := GetLatestBid(bidId)
	if err != nil {
		return pgconn.CommandTag{}, err, currentBid.Version
	}
	newVersion := currentBid.Version + 1
	tag, err := database.DataSource.Exec(context.Background(),
		`INSERT INTO bids (tender_id, description, status, version, created_at, updated_at, bid_id, username, name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		currentBid.TenderID, update.Description, currentBid.Status, newVersion, time.Now(), time.Now(), currentBid.BidID, currentBid.UserName, update.Name)
	return tag, err, newVersion
}

func RollbackBid(bidId int, version int) (*models.Bid, error) {
	rollbackBid, err := GetBidByVersion(bidId, version)
	if err != nil {
		return rollbackBid, err
	}
	_, err, latestVersion := UpdateBid(rollbackBid, bidId)
	if err != nil {
		return rollbackBid, err
	}

	var newBid models.Bid
	err = database.DataSource.QueryRow(context.Background(), `SELECT * FROM bids WHERE bid_id=$1 AND version=$2`, bidId, latestVersion).Scan(
		&newBid.ID, &newBid.TenderID, &newBid.Description, &newBid.Status, &newBid.Version, &newBid.CreatedAt, &newBid.UpdatedAt, &newBid.BidID, &newBid.UserName, &newBid.Name)

	return &newBid, err
}

func GetBidByVersion(bidId int, version int) (*models.Bid, error) {
	var oldBid models.Bid
	err := database.DataSource.QueryRow(context.Background(), `SELECT * FROM bids WHERE bid_id=$1 AND version=$2`, bidId, version).Scan(&oldBid.ID, &oldBid.TenderID, &oldBid.Description, &oldBid.Status, &oldBid.Version, &oldBid.CreatedAt, &oldBid.UpdatedAt, &oldBid.BidID, &oldBid.UserName, &oldBid.Name)

	if err != nil {
		log.Printf("Error retrieving bid with version %d: %v", version, err)
		return &models.Bid{}, err
	}
	return &oldBid, err
}

func GetLatestBid(bidId int) (*models.Bid, error) {
	var currentBid models.Bid
	err := database.DataSource.QueryRow(context.Background(), `SELECT * FROM bids WHERE bid_id=$1 AND version=(SELECT MAX(version) FROM bids WHERE bid_id=$1)`, bidId).Scan(&currentBid.ID, &currentBid.TenderID, &currentBid.Description, &currentBid.Status, &currentBid.Version, &currentBid.CreatedAt, &currentBid.UpdatedAt, &currentBid.BidID, &currentBid.UserName, &currentBid.Name)
	return &currentBid, err
}

func MakeDecision(decision string, bidId int) error {
	_, err := database.DataSource.Exec(context.Background(), `UPDATE bids SET status=$1 WHERE bid_id=$2`, decision, bidId)
	if err != nil {
		return err
	}

	if decision == "Approved" {
		var tenderID int
		err = database.DataSource.QueryRow(context.Background(), `SELECT tender_id FROM bids WHERE bid_id=$1`, bidId).Scan(&tenderID)
		if err != nil {
			return err
		}

		_, err = database.DataSource.Exec(context.Background(), `UPDATE tenders SET status='CLOSED' WHERE tender_id=$1`, tenderID)
		if err != nil {

			return err
		}
	}
	return err
}

func GetBidsByTenderId(tenderId int) (pgx.Rows, error) {
	return database.DataSource.Query(context.Background(), `SELECT * FROM bids WHERE tender_id=$1`, tenderId)
}
