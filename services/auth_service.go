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

func Login(loginRequest LoginRequest) (*models.UserData, *util.TokenPair, error) {
	cardData, err := util.ParseScannerString(loginRequest.Scanner_Message)
	tokenPair := &util.TokenPair{}
	if err != nil {
        return nil, tokenPair, fmt.Errorf("error parsing card data: %v", err)
    }

	var userData models.UserData
	err = database.DB.QueryRow("SELECT id, username, has_training, admin, has_executive_access, is_egn_lab, ban_time_end, weekly_minutes FROM users WHERE id = ?", cardData.Id).Scan(
		&userData.Id,
		&userData.Username,
		&userData.Trained,
		&userData.Admin,
		&userData.Has_Executive_Access,
		&userData.Is_Egn_Lab,
		&userData.Ban_Time_End,
		&userData.Weekly_Minutes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
            return nil, tokenPair, ErrorUserNotFound
		}
		return nil, tokenPair, fmt.Errorf("database error: %v", err)
	}

	if !userData.Trained {
        return nil, tokenPair, ErrorNotTrained
    }

	token, err := util.GenerateTokenPair(userData.Id, userData.Admin)
	if err != nil {
        return nil, tokenPair , fmt.Errorf("error generating token: %v", err)
    }
	
	return &userData, token, nil
}

