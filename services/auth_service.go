package services

import (
	"database/sql"
	"errors"
	"fmt"
	"gin-api/database"
	"gin-api/models"
	"gin-api/util"
)

type LoginRequest struct {
	Scanner_Message string `json:"scanner_message"`
}

var (
	ErrorUserNotFound = errors.New("user not found")
	ErrorNotTrained   = errors.New("user not trained")
)

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
            return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	if !userData.Trained {
        return nil, ErrorNotTrained
    }
	
	return &userData, nil
}

