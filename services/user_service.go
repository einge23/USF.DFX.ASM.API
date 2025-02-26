package services

import (
	"database/sql"
	"fmt"
	"gin-api/database"
	"gin-api/models"
	"gin-api/util"
)

type CreateUserRequest struct {
	Scanner_Message string `json:"scanner_message"`
	Trained         bool
	Admin           bool
}

func CreateUser(createUserRequest CreateUserRequest) (bool, error) {

	cardData, err := util.ParseScannerString(createUserRequest.Scanner_Message)
	if err != nil {
		return true, fmt.Errorf("error parsing card data: %v", err)
	}

	//check for existing user
	var existingID int
	querySQL := `SELECT id FROM users WHERE username = ?`
	err = database.DB.QueryRow(querySQL, cardData.Username).Scan(&existingID)

	//problem querying
	if err != nil && err != sql.ErrNoRows {
		return true, fmt.Errorf("could not query user: %v", err)
	}

	//user already exists
	if existingID != 0 {
		return true, fmt.Errorf("user with username '%s' already exists", cardData.Username)
	}

	//add user
	insertSQL := `INSERT INTO users (id, username, has_training, admin) VALUES (?, ?, ?, ?)`
	_, err = database.DB.Exec(insertSQL, cardData.Id, cardData.Username, createUserRequest.Trained, createUserRequest.Admin)
	if err != nil {
		return true, fmt.Errorf("could not add user: %v", err)
	}
	return false, nil
}

type SetUserTrainedRequest struct {
	UserId int `json:"user_id"`
}

func SetUserTrained(setUserTrainedRequest SetUserTrainedRequest) (bool, error) {

	//get trained status for the user from db
	var trainedStatus bool
	fmt.Println(setUserTrainedRequest.UserId)
	querySQL := `SELECT has_training FROM users WHERE id = ?`
	err := database.DB.QueryRow(querySQL, setUserTrainedRequest.UserId).Scan(&trainedStatus)
	if err != nil {
		return true, fmt.Errorf("error getting user from db: %v", err)
	}

	//toggle the trainedStatus
	newTrainedStatus := !trainedStatus

	//update the user's training status in the database
	updateSQL := `UPDATE users SET has_training = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL, newTrainedStatus, setUserTrainedRequest.UserId)
	if err != nil {
		return true, fmt.Errorf("error updating user training status: %v", err)
	}

	return false, nil
}

func GetUserReservations(userId int) ([]models.ReservationDTO, error) {
	var reservations []models.ReservationDTO
	querySQL := `SELECT id, userId, time_reserved, time_complete, printerid, is_active FROM reservations WHERE userId = ? ORDER BY time_reserved DESC`
	rows, err := database.DB.Query(querySQL, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting reservations from db: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reservation models.ReservationDTO
		err := rows.Scan(&reservation.Id, &reservation.UserId, &reservation.TimeReserved, &reservation.TimeComplete, &reservation.PrinterId, &reservation.IsActive)
		if err != nil {
			return nil, fmt.Errorf("error scanning reservation: %v", err)
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func GetActiveUserReservations(userId int) ([]models.ReservationDTO, error) {
	var reservations []models.ReservationDTO
	querySQL := `SELECT id, userId, time_reserved, time_complete, printerid, is_active FROM reservations WHERE userId = ? AND is_active = 1 ORDER BY time_reserved DESC`
	rows, err := database.DB.Query(querySQL, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting reservations from db: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reservation models.ReservationDTO
		err := rows.Scan(&reservation.Id, &reservation.UserId, &reservation.TimeReserved, &reservation.TimeComplete, &reservation.PrinterId, &reservation.IsActive)
		if err != nil {
			return nil, fmt.Errorf("error scanning reservation: %v", err)
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func GetUserById(userID int) (*models.UserData, error) {
	var user models.UserData
	querySQL := `SELECT id, username, has_training, admin, has_executive_access FROM users WHERE id = ?`
	err := database.DB.QueryRow(querySQL, userID).Scan(&user.Id, &user.Username, &user.Trained, &user.Admin, &user.Has_Executive_Access)
	if err != nil {
		return nil, fmt.Errorf("error getting user from db: %v", err)
	}
	return &user, nil
}

type SetUserExecutiveAccessRequest struct {
	UserId int `json:"user_id"`
}

func SetUserExecutiveAccess(setUserExecutiveAccessRequest SetUserExecutiveAccessRequest) (bool, error) {

	var currentExecutiveAccess bool

	querySQL := `SELECT has_executive_access FROM users WHERE id = ?`
	err := database.DB.QueryRow(querySQL, setUserExecutiveAccessRequest.UserId).Scan(&currentExecutiveAccess)
	if err != nil {
		return true, fmt.Errorf("error getting user executive access from db: %v", err)
	}

	newExecutiveAccess := !currentExecutiveAccess

	updateSQL := `UPDATE users SET has_executive_access = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL, newExecutiveAccess, setUserExecutiveAccessRequest.UserId)
	if err != nil {
		return true, fmt.Errorf("error updating user executive access: %v", err)
	}

	return false, nil //return 0
}
