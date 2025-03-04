package services

import (
	"database/sql"
	"fmt"
	"gin-api/database"
	"gin-api/models"
	"gin-api/util"
	"time"
)

type CreateUserRequest struct {
	Scanner_Message string `json:"scanner_message"`
	Trained         bool   `json:"trained"`
	Admin           bool   `json:"admin"`
}

func CreateUser(createUserRequest CreateUserRequest) (bool, error) {

	cardData, err := util.ParseScannerString(createUserRequest.Scanner_Message)
	if err != nil {
		return false, fmt.Errorf("error parsing card data: %v", err)
	}

	//check for existing user
	var existingID int
	querySQL := `SELECT id FROM users WHERE username = ?`
	err = database.DB.QueryRow(querySQL, cardData.Username).Scan(&existingID)

	//problem querying
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("could not query user: %v", err)
	}

	//user already exists
	if existingID != 0 {
		return false, fmt.Errorf("user with username '%s' already exists", cardData.Username)
	}

	//add user
	insertSQL := `INSERT INTO users (id, username, has_training, admin) VALUES (?, ?, ?, ?)`
	_, err = database.DB.Exec(insertSQL, cardData.Id, cardData.Username, createUserRequest.Trained, createUserRequest.Admin)
	if err != nil {
		return false, fmt.Errorf("could not add user: %v", err)
	}
	return true, nil
}

func SetUserTrained(userId int) error {

	//get trained status for the user from db
	var trainedStatus bool
	querySQL := `SELECT has_training FROM users WHERE id = ?`
	err := database.DB.QueryRow(querySQL, userId).Scan(&trainedStatus)
	if err != nil {
		return fmt.Errorf("error getting user from db: %v", err)
	}

	//toggle the trainedStatus
	newTrainedStatus := !trainedStatus

	//update the user's training status in the database
	updateSQL := `UPDATE users SET has_training = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL, newTrainedStatus, userId)
	if err != nil {
		return fmt.Errorf("error updating user training status: %v", err)
	}

	return nil
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

func SetUserExecutiveAccess(userId int) error {

	var currentExecutiveAccess bool

	querySQL := `SELECT has_executive_access FROM users WHERE id = ?`
	err := database.DB.QueryRow(querySQL, userId).Scan(&currentExecutiveAccess)
	if err != nil {
		return fmt.Errorf("error getting user executive access from db: %v", err)
	}

	newExecutiveAccess := !currentExecutiveAccess

	updateSQL := `UPDATE users SET has_executive_access = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL, newExecutiveAccess, userId)
	if err != nil {
		return fmt.Errorf("error updating user executive access: %v", err)
	}

	return nil
}

type AddUserWeeklyMinutesRequest struct {
	Minutes int `json:"minutes"`
}

func AddUserWeeklyMinutes(id int, request AddUserWeeklyMinutesRequest) error {
	var currentWeeklyMinutes int

	querySQL := `SELECT weekly_minutes FROM users WHERE id = ?`
	err := database.DB.QueryRow(querySQL, id).Scan(&currentWeeklyMinutes)
	if err != nil {
		return fmt.Errorf("error getting user weekly minutes from db: %v", err)
	}

	newWeeklyMinutes := currentWeeklyMinutes + request.Minutes

	updateSQL := `UPDATE users SET weekly_minutes = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL, newWeeklyMinutes, id)
	if err != nil {
		return fmt.Errorf("error adding minutes to user: %v", err)
	}

	return nil
}

type SetUserBanTimeRequest struct {
	BanTime int `json:"ban_time"`
}

// Set the user's ban time. Grab id from path, get time in hours from json body. A -1 in json body means set the
// user's ban_time_end back to NULL in the database. If there is no ban time, the ban becomes current time + hours
// requested. If there is a current ban, the ban becomes current ban date and time + requested hours, extending the ban.
func SetUserBanTime(id int, request SetUserBanTimeRequest) error {

	if request.BanTime == -1 { //if passing in -1, set ban_time_end back to NULL in db
		updateSQL := `UPDATE users SET ban_time_end = ? WHERE id = ?`
		_, err := database.DB.Exec(updateSQL, nil, id)
		if err != nil {
			return fmt.Errorf("error setting ban time to NULL for user: %v", err)
		}
		return nil
	}
	var currentBanTimeEnd *time.Time

	querySQL := `SELECT ban_time_end FROM users WHERE id = ?`
	err := database.DB.QueryRow(querySQL, id).Scan(&currentBanTimeEnd)
	if err != nil {
		return fmt.Errorf("error getting user ban time from db: %v", err)
	}

	var newBanTimeEnd time.Time

	if currentBanTimeEnd == nil { //if null, set to current time + requested ban time
		newBanTimeEnd = time.Now().Add(time.Duration(request.BanTime) * time.Hour)
	} else { //if not null, set to existing ban time end + requested ban time
		newBanTimeEnd = currentBanTimeEnd.Add(time.Duration(request.BanTime) * time.Hour)
	}

	updateSQL := `UPDATE users SET ban_time_end = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL, newBanTimeEnd, id)
	if err != nil {
		return fmt.Errorf("error adding ban time to user: %v", err)
	}

	return nil
}
