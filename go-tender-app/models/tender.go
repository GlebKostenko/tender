package models

import "time"

type TenderStatus string

const (
	Created   TenderStatus = "Created"
	Published TenderStatus = "Published"
	Closed    TenderStatus = "Closed"
)

type Tender struct {
	ID             int          `json:"id"`
	TenderID       int          `json:"tenderId"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	OrganizationID int          `json:"organizationId"`
	Status         TenderStatus `json:"status"`
	UserName       string       `json:"creatorUsername"`
	Version        int          `json:"version"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}
