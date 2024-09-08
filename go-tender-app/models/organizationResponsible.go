package models

type OrganizationResponsible struct {
	ID             int `json:"id"`
	OrganizationId int `json:"organizationId"`
	UserId         int `json:"userId"`
}
