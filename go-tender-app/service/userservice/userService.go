package userservice

import (
	"context"
	"go-tender-app/database"
	"log"
)

func UserIsUserResponsible(userName *string, organizationId int) (bool, error) {
	userId, err := GetUserIdByName(userName)
	if err != nil {
		return false, err
	}
	var isResponsible bool
	err = database.DataSource.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM organization_responsible WHERE user_id=$1 AND organization_id=$2)`,
		userId, organizationId).
		Scan(&isResponsible)
	if err != nil {
		return false, err
	}
	return isResponsible, err
}

func GetUserIdByName(userName *string) (int, error) {
	var userId int
	err := database.DataSource.QueryRow(context.Background(),
		`SELECT id FROM employee WHERE username=$1`, userName).
		Scan(&userId)
	if err != nil {
		log.Printf("Failed to find user with name %v", userName)
	}
	return userId, err
}
