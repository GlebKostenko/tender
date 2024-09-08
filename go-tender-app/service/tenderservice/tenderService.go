package tenderservice

import (
	"context"
	"go-tender-app/database"
	"go-tender-app/models"
	"go-tender-app/service/userservice"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

func GetAllTenders() (pgx.Rows, error) {
	return database.DataSource.Query(context.Background(), "SELECT * FROM tenders")
}

func Validate(rows *pgx.Rows) (*[]models.Tender, error) {
	var err error
	var tenders []models.Tender
	for (*rows).Next() {
		var tender models.Tender
		if err = (*rows).Scan(&tender.ID, &tender.OrganizationID, &tender.Name, &tender.UserName, &tender.Description, &tender.Status, &tender.Version, &tender.CreatedAt, &tender.UpdatedAt, &tender.TenderID); err != nil {
			return &tenders, err
		}
		tenders = append(tenders, tender)
	}
	return &tenders, err
}

func Create(tender *models.Tender) error {
	isResponsible, err := userservice.UserIsUserResponsible(&tender.UserName, tender.OrganizationID)
	if err != nil || !isResponsible {
		log.Printf("Unauthorized to create a tender by user: %v", tender.UserName)
		return err
	}

	var maxTenderID int
	err = database.DataSource.QueryRow(context.Background(), "SELECT COALESCE(MAX(tender_id), 0) FROM tenders").Scan(&maxTenderID)
	if err != nil {
		return err
	}
	tenderId := maxTenderID + 1
	return database.DataSource.QueryRow(context.Background(),
		`INSERT INTO tenders (name, tender_id, description, organization_id, status, version, username) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		tender.Name, tenderId, tender.Description, tender.OrganizationID, models.Created, 1, tender.UserName).
		Scan(&tender.ID, &tender.CreatedAt, &tender.UpdatedAt)
}

func GetByUsername(userName *string) (pgx.Rows, error) {
	userId, err := userservice.GetUserIdByName(userName)
	if err != nil {
		return nil, err
	}
	return database.DataSource.Query(
		context.Background(),
		`SELECT * FROM tenders WHERE organization_id IN (SELECT organization_id FROM organization_responsible WHERE user_id=$1)`,
		userId)

}

func GetStatus(tenderID int) (*models.TenderStatus, error) {
	var status models.TenderStatus
	err := database.DataSource.QueryRow(context.Background(), `SELECT status FROM tenders WHERE tender_id=$1 AND version=(SELECT MAX(version) FROM tenders WHERE tender_id=$1)`, tenderID).Scan(&status)
	return &status, err
}

func UpdateTender(tender *models.Tender, tenderId int) (pgconn.CommandTag, error, int) {
	curTender, err := GetLatestTender(tenderId)
	if err != nil {
		return pgconn.CommandTag{}, err, curTender.Version
	}
	newVersion := curTender.Version + 1

	tag, err := database.DataSource.Exec(context.Background(),
		`INSERT INTO tenders (name, description, status, version, created_at, updated_at, username, organization_id, tender_id)
	 VALUES ($1, $2, $3, $4, $5, NOW(), $6, $7, $8)`,
		tender.Name, tender.Description, curTender.Status, newVersion, tender.CreatedAt, curTender.UserName, curTender.OrganizationID, curTender.TenderID)
	return tag, err, newVersion
}

func RollbackTender(tenderId int, version int) (*models.Tender, error) {
	rollbackTender, err := GetTenderByVersion(tenderId, version)
	if err != nil {
		return &models.Tender{}, err
	}
	_, err, latestVersion := UpdateTender(rollbackTender, tenderId)
	if err != nil {
		return &models.Tender{}, err
	}

	var newTender models.Tender
	err = database.DataSource.QueryRow(context.Background(),
		`SELECT id, name, username, description, status, version, created_at, updated_at 
         FROM tenders 
         WHERE tender_id=$1 AND version=$2`, tenderId, latestVersion).
		Scan(&newTender.ID, &newTender.Name, &newTender.UserName, &newTender.Description, &newTender.Status, &newTender.Version, &newTender.CreatedAt, &newTender.UpdatedAt)
	return &newTender, err

}

func GetTenderByVersion(tenderId int, version int) (*models.Tender, error) {
	var rollbackTender models.Tender
	err := database.DataSource.QueryRow(context.Background(),
		`SELECT id, name, username, description, status, version, created_at, updated_at 
         FROM tenders 
         WHERE tender_id=$1 AND version=$2`, tenderId, version).
		Scan(&rollbackTender.ID, &rollbackTender.Name, &rollbackTender.UserName, &rollbackTender.Description, &rollbackTender.Status, &rollbackTender.Version, &rollbackTender.CreatedAt, &rollbackTender.UpdatedAt)
	if err != nil {
		log.Printf("Error retrieving tender with version %d: %v", version, err)
		return &models.Tender{}, err
	}
	return &rollbackTender, err
}

func GetLatestTender(tenderId int) (*models.Tender, error) {
	var currentTender models.Tender
	err := database.DataSource.QueryRow(context.Background(), `SELECT * FROM tenders WHERE tender_id=$1 AND version=(SELECT MAX(version) FROM tenders WHERE tender_id=$1)`, tenderId).Scan(&currentTender.ID, &currentTender.OrganizationID, &currentTender.Name, &currentTender.UserName, &currentTender.Description, &currentTender.Status, &currentTender.Version, &currentTender.CreatedAt, &currentTender.UpdatedAt, &currentTender.TenderID)
	return &currentTender, err
}

func Publish(tenderId int) (pgconn.CommandTag, error) {
	return database.DataSource.Exec(context.Background(), `UPDATE tenders SET status='PUBLISHED' WHERE tender_id=$1`, tenderId)
}

func Close(tenderId int) (pgconn.CommandTag, error) {
	return database.DataSource.Exec(context.Background(), `UPDATE tenders SET status='CLOSED' WHERE tender_id=$1`, tenderId)
}
