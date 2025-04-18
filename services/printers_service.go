package services

import (
	"database/sql"
	"fmt"
	"gin-api/database"
	"gin-api/models"
	"gin-api/util"
	"log"
	"time"
)

// return all printers by rack for specific system (EGN or General), as serialized JSON
func GetPrinters(isEgnLab bool) ([]models.Printer, error) {
	// Build query based on isEgnLab parameter
	query := "SELECT id, name, color, rack, in_use, last_reserved_by, is_executive FROM printers WHERE is_egn_printer = ? order by rack asc"

	// Execute query with appropriate parameter
	rows, err := database.DB.Query(query, isEgnLab)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var printers []models.Printer
	for rows.Next() {
		var p models.Printer
		var lastReservedBy sql.NullString
		if err := rows.Scan(&p.Id, &p.Name, &p.Color, &p.Rack, &p.In_Use, &lastReservedBy, &p.Is_Executive); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		if lastReservedBy.Valid {
			p.Last_Reserved_By = lastReservedBy.String
		} else {
			p.Last_Reserved_By = ""
		}
		printers = append(printers, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return printers, nil
}

// given a printer object, add a printer with those attributes. ID correlates to physical plug.
func AddPrinter(request models.Printer) (bool, error) {

	var id int
	row := database.DB.QueryRow("SELECT id FROM printers WHERE id = ?", request.Id)
	err := row.Scan(&id)

	//No row exists, proceed with add
	if err == sql.ErrNoRows {
		insertSQL := `INSERT INTO printers (id, name, color, rack, in_use, last_reserved_by, is_executive, is_egn_printer) values (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = database.DB.Exec(insertSQL,
			request.Id,
			request.Name,
			request.Color,
			request.Rack,
			request.In_Use,
			nil,
			request.Is_Executive,
			request.Is_Egn_Printer)
		if err != nil {
			return false, fmt.Errorf("error inserting new printer to DB: %v", err)
		}
		return true, nil

		//some other error has happened
	} else if err != nil {
		return false, err
	}
	//Row exists
	return false, fmt.Errorf("printer with specified ID already exists")
}

type UpdatePrinterRequest struct {
	Name         string `json:"name"`
	Color        string `json:"color"`
	Rack         int    `json:"rack"`
	IsExecutive  bool   `json:"is_executive"`
	IsEgnPrinter bool   `json:"is_egn_printer"`
}

// given printer id and name, color, rack, is_executive, is_egn, change the values of that printer
// to match the attributes passed in.
func UpdatePrinter(id int, request UpdatePrinterRequest) (bool, error) {

	row := database.DB.QueryRow("SELECT id FROM printers WHERE id = ?", id)
	err := row.Scan(&id)

	//No row exists
	if err == sql.ErrNoRows {
		return false, fmt.Errorf("printer does not exist")

		//some other error has happened
	} else if err != nil {
		return false, err
	}
	//printer exists
	updateSQL := `UPDATE printers SET name = ?, color = ?, rack = ?, is_executive = ?, is_egn_printer = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL,
		request.Name,
		request.Color,
		request.Rack,
		request.IsExecutive,
		request.IsEgnPrinter,
		id)
	if err != nil {
		return false, fmt.Errorf("error updating printer in DB: %v", err)
	}
	return true, nil
}

type ReservePrinterRequest struct {
	PrinterId int `json:"printer_id"`
	UserId    int `json:"user_id"`
	TimeMins  int `json:"time_mins"`
}

var (
	manager = &models.ReservationManager{
		Reservations: make(map[int]*models.Reservation),
	}
)

// given a printerId, userId, and time in minutes, reserve that printer for the user
// and for that many minutes. Also add a timed event to complete the reservation after
// the time in minutes has passed.
func ReservePrinter(printerId int, userId int, timeMins int) (bool, error) {
	var user models.UserData
	if err := database.DB.QueryRow("SELECT username FROM users WHERE id = ?", userId).Scan(
		&user.Username); err != nil {
		return false, fmt.Errorf("failed to get username: %v", err)
	}

	var printer models.Printer
	var lastReservedBy sql.NullString
	if err := database.DB.QueryRow("SELECT id, name, color, rack, in_use, last_reserved_by, is_executive, is_egn_printer FROM printers WHERE id = ?", printerId).Scan(
		&printer.Id,
		&printer.Name,
		&printer.Color,
		&printer.Rack,
		&printer.In_Use,
		&lastReservedBy,
		&printer.Is_Executive,
		&printer.Is_Egn_Printer); err != nil {
		return false, fmt.Errorf("failed to get printer: %v", err)
	}

	if lastReservedBy.Valid {
		printer.Last_Reserved_By = lastReservedBy.String
	} else {
		printer.Last_Reserved_By = ""
	}

	if printer.In_Use {
		return false, fmt.Errorf("printer is already in use")
	}

	// Check if user already has all of his active reservations
	var activeReservationCount int
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM reservations WHERE userid = ? AND is_active = TRUE", userId).Scan(&activeReservationCount); err != nil {
		return false, fmt.Errorf("failed to check active reservations: %v", err)
	}

	// Check if the active reservation count is larger than the limit and return false if it is
	limit := util.Settings.PrinterSettings.MaxActiveReservations // Set limit to the amount of active reservations that the administrator passes
	if activeReservationCount >= limit {
		return false, fmt.Errorf("maximum of active reservations per user allowed is %d", limit)
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Use defer with a named error to handle rollback/commit
	var txErr error
	defer func() {
		if txErr != nil {
			tx.Rollback()
		}
	}()

	// Set printer as 'in use' in the database (within transaction)
	result, err := tx.Exec(
		"UPDATE printers SET in_use = TRUE, last_reserved_by = ? WHERE id = ?",
		user.Username,
		printerId,
	)
	if err != nil {
		txErr = err
		return false, fmt.Errorf("failed to update printer: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		txErr = err
		return false, fmt.Errorf("failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		txErr = fmt.Errorf("no printer found with id: %d", printerId)
		return false, txErr
	}

	time_reserved := time.Now()
	time_complete := time.Now().Add(time.Duration(timeMins) * time.Minute)

	// Create reservation and add it to reservations table as an entry (within transaction)
	result, err = tx.Exec(
		"INSERT INTO reservations (printerid, userid, time_reserved, time_complete, is_active, is_egn_reservation) values (?, ?, ?, ?, ?, ?)",
		printerId,
		userId,
		time_reserved,
		time_complete,
		true,
		printer.Is_Egn_Printer)
	if err != nil {
		txErr = err
		return false, fmt.Errorf("failed to insert reservation: %v", err)
	}

	reservationId, err := result.LastInsertId()
	if err != nil {
		txErr = err
		return false, fmt.Errorf("failed to get reservation id: %v", err)
	}

	// Get user's weeklyMinutes (within transaction)
	var currentWeeklyMinutes int
	err = tx.QueryRow("SELECT weekly_minutes FROM users WHERE id = ?", userId).Scan(&currentWeeklyMinutes)
	if err != nil {
		txErr = err
		return false, fmt.Errorf("error getting user weekly minutes from db: %v", err)
	}

	// Subtract reservation's duration from weeklyMinutes (within transaction)
	newWeeklyMinutes := currentWeeklyMinutes - timeMins
	_, err = tx.Exec("UPDATE users SET weekly_minutes = ? WHERE id = ?", newWeeklyMinutes, userId)
	if err != nil {
		txErr = err
		return false, fmt.Errorf("error subtracting minutes from user: %v", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Only turn on printer after successful transaction
	_, err = util.TurnOnPrinter(printerId)
	if err != nil {
		// Transaction was successful but printer failed to turn on
		// We should try to undo our changes
		undoErr := undoReservation(printerId, int(reservationId), userId, currentWeeklyMinutes)
		if undoErr != nil {
			log.Printf("failed to undo reservation after printer turn on error: %v", undoErr)
		}
		return false, fmt.Errorf("error turning on printer: %v", err)
	}

	// Set up timer to complete/end the reservation
	timer := time.NewTimer(time.Duration(timeMins) * time.Minute)
	manager.Mutex.Lock()
	manager.Reservations[int(reservationId)] = &models.Reservation{
		Id:                 int(reservationId),
		PrinterId:          printerId,
		UserId:             userId,
		Time_Reserved:      time_reserved,
		Time_Complete:      time_complete,
		Is_Active:          true,
		Is_Egn_Reservation: printer.Is_Egn_Printer,
		Timer:              timer,
	}
	manager.Mutex.Unlock()

	go func() {
		<-timer.C
		CompleteReservation(printerId, int(reservationId))
	}()

	return true, nil
}

// Helper function to undo a reservation if printer fails to turn on
func undoReservation(printerId, reservationId, userId, originalWeeklyMinutes int) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin undo transaction: %v", err)
	}

	// Set printer back to not in use
	_, err = tx.Exec("UPDATE printers SET in_use = FALSE WHERE id = ?", printerId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to undo printer status: %v", err)
	}

	// Set reservation to inactive
	_, err = tx.Exec("UPDATE reservations SET is_active = FALSE WHERE id = ?", reservationId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to undo reservation: %v", err)
	}

	// Restore user's weekly minutes
	_, err = tx.Exec("UPDATE users SET weekly_minutes = ? WHERE id = ?", originalWeeklyMinutes, userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to restore user minutes: %v", err)
	}

	return tx.Commit()
}

// turn off the printer, set the relevant printer as not in use, set the reservation to no longer be active
func CompleteReservation(printerId, reservationId int) {

	//Turn off the printer
	_, err := util.TurnOffPrinter(printerId)
	if err != nil {
		log.Printf("failed to turn off printer: %v", err)
	}

	//Set as not in_use
	_, err = database.DB.Exec(
		"UPDATE printers SET in_use = FALSE WHERE id = ?",
		printerId,
	)
	if err != nil {
		log.Printf("failed to update printer: %v", err)
	}
	//Set the reservation as inactive
	_, err = database.DB.Exec(
		"UPDATE reservations SET is_active = FALSE WHERE id = ?",
		reservationId,
	)
	if err != nil {
		log.Printf("failed to update reservation: %v", err)
	}

	manager.Mutex.Lock()
	delete(manager.Reservations, reservationId)
	manager.Mutex.Unlock()
}

// Given a printerId, toggle its is_executive bool in the printers table
func SetPrinterExecutive(id int) error {

	var currentExecutiveness bool

	querySQL := `SELECT is_executive FROM printers WHERE id = ?`
	err := database.DB.QueryRow(querySQL, id).Scan(&currentExecutiveness)
	if err != nil {
		return fmt.Errorf("error getting printer executiveness from db: %v", err)
	}

	newExecutiveness := !currentExecutiveness

	updateSQL := `UPDATE printers SET is_executive = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL, newExecutiveness, id)
	if err != nil {
		return fmt.Errorf("error updating printer executiveness: %v", err)
	}

	return nil
}
