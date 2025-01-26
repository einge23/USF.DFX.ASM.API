package services

import (
	"database/sql"
	"fmt"
	"gin-api/database"
	"gin-api/models"
	"gin-api/util"
)

type LoginRequest struct {
	Scanner_Message string `json:"scanner_message"`
}

func Login(loginRequest LoginRequest) (*models.UserData, error) {
	cardData, err := util.ParseScannerString(loginRequest.Scanner_Message)
	if err != nil {
        return nil, fmt.Errorf("error parsing card data: %v", err)
    }

	var userData models.UserData
	err = database.DB.QueryRow("SELECT id, username, has_training, admin FROM users WHERE id = ?", cardData.Id).Scan(
		&userData.Id,
		&userData.Username,
		&userData.Trained,
		&userData.Admin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	
	return &userData, nil
}

