package models

import "time"

type BidStatus string
type AuthorT string

const (
	BidCreated   BidStatus = "Created"
	BidPublished BidStatus = "Published"
	BidCanceled  BidStatus = "Closed"
	BidApproved  BidStatus = "Approved"
	BidRejected  BidStatus = "Rejected"
)

const (
	AuthorOrganization AuthorT = "Organization"
	AuthorUser         AuthorT = "User"
)

type Bid struct {
	ID          int       `json:"id"`
	BidID       int       `json:"bidId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      BidStatus `json:"status"`
	TenderID    int       `json:"tenderId"`
	UserName    string    `json:"creatorUsername"`
	AuthorType  AuthorT   `json:"authorType"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
