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
	Egn_Lab         bool `json:"egn_lab"`
}

//Given a card scanner raw input, trained bool, and admin bool, create a user and add it to user table
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
	insertSQL := `INSERT INTO users (id, username, has_training, admin, is_egn_lab) VALUES (?, ?, ?, ?, ?)`
	_, err = database.DB.Exec(insertSQL, cardData.Id, cardData.Username, createUserRequest.Trained, createUserRequest.Admin, createUserRequest.Egn_Lab)
	if err != nil {
		return false, fmt.Errorf("could not add user: %v", err)
	}
	return true, nil
}

//Given a userId, toggle that user's has_training bool in the users table
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

//given a userId, return all reservation info of that user. Returns both active and inactive reservations.
func GetUserReservations(userId int) ([]models.ReservationDTO, error) {
	var reservations []models.ReservationDTO
	querySQL := `
		SELECT 
			r.id, r.userId, u.username, r.time_reserved, r.time_complete, 
			r.printerid, p.name AS printer_name, r.is_active, r.is_egn_reservation 
		FROM reservations r
		JOIN users u ON r.userId = u.id
		JOIN printers p ON r.printerid = p.id
		WHERE r.userId = ? 
		ORDER BY r.time_reserved DESC
	`
	rows, err := database.DB.Query(querySQL, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting reservations from db: %v", err)
	}
	defer rows.Close()

	for rows.Next() { //for all returns rows
		var reservation models.ReservationDTO //create a reservation object
		//scan in all reservation info
		err := rows.Scan(
			&reservation.Id, &reservation.UserId, &reservation.Username, &reservation.Time_Reserved,
			&reservation.Time_Complete, &reservation.PrinterId, &reservation.PrinterName,
			&reservation.Is_Active, &reservation.Is_Egn_Reservation,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning reservation: %v", err)
		}
		//add the reservation object to a list
		reservations = append(reservations, reservation)
	}
	//return the whole list of reservations
	return reservations, nil
}

//given a userId, return all reservation info for active reservations of that user
func GetActiveUserReservations(userId int) ([]models.ReservationDTO, error) {
	var reservations []models.ReservationDTO
	querySQL := `
		SELECT 
			r.id, r.userId, u.username, r.time_reserved, r.time_complete, 
			r.printerid, p.name AS printer_name, r.is_active, r.is_egn_reservation 
		FROM reservations r
		JOIN users u ON r.userId = u.id
		JOIN printers p ON r.printerid = p.id
		WHERE r.userId = ? AND r.is_active = 1 
		ORDER BY r.time_reserved DESC
	`
	rows, err := database.DB.Query(querySQL, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting reservations from db: %v", err)
	}
	defer rows.Close()

	for rows.Next() { //for all returned rows
		var reservation models.ReservationDTO //create a reservation object
		//scan in all active reservation info
		err := rows.Scan(
			&reservation.Id, &reservation.UserId, &reservation.Username, &reservation.Time_Reserved,
			&reservation.Time_Complete, &reservation.PrinterId, &reservation.PrinterName,
			&reservation.Is_Active, &reservation.Is_Egn_Reservation,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning reservation: %v", err)
		}
		//add the reservation object to a list
		reservations = append(reservations, reservation)
	}
	//return the whole list of reservations
	return reservations, nil
}

//given a userId, return a user object with all user data
func GetUserById(userID int) (*models.UserData, error) {
	var user models.UserData
	querySQL := `SELECT id, username, has_training, admin, has_executive_access, is_egn_lab, ban_time_end, weekly_minutes FROM users WHERE id = ?`
	err := database.DB.QueryRow(querySQL, userID).Scan(&user.Id, &user.Username, &user.Trained, &user.Admin, &user.Has_Executive_Access, &user.Is_Egn_Lab, &user.Ban_Time_End, &user.Weekly_Minutes)
	if err != nil {
		return nil, fmt.Errorf("error getting user from db: %v", err)
	}
	return &user, nil
}

//given a userId, toggle the user's has_executive_access bool in the users table
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

//given a userId and a number of minutes, add those minutes to the user's weekly_minutes
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

//given a userId and a number of hours, add that number of hours to the user's ban. If the user is not banned, they
//are banned until (now plus the requested hours). If they are banned, add the requested hours to their existing ban time
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
