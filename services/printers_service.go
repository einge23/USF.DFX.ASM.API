package services

import (
	"database/sql"
	"errors"
	"fmt"
	"gin-api/database"
	"gin-api/models"
	"gin-api/util"
	"log"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3" // Import the sqlite3 driver
)

// return all printers by rack as serialized JSON
func GetPrinters() ([]models.Printer, error) {
	// Build query
	query := "SELECT id, name, color, rack, rack_position, in_use, last_reserved_by, is_executive FROM printers order by rack asc, rack_position asc"

	// Execute query with appropriate parameter
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var printers []models.Printer
	for rows.Next() {
		var p models.Printer
		var lastReservedBy sql.NullString
		// Scan rack_position
		if err := rows.Scan(&p.Id, &p.Name, &p.Color, &p.Rack, &p.Rack_Position, &p.In_Use, &lastReservedBy, &p.Is_Executive); err != nil {
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

// given a printer object, add a printer with those attributes. ID correlates to physical plug (1-28).
// Rack position is automatically calculated as the next available position in the specified rack.
func AddPrinter(request models.Printer) (bool, error) {
	// --- BEGINNING OF CHANGES ---
	// Validate Printer ID range
	if request.Id < 1 || request.Id > 28 {
		return false, fmt.Errorf("invalid printer ID: %d. ID must be between 1 and 28", request.Id)
	}

	// Check total printer count
	var printerCount int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM printers").Scan(&printerCount)
	if err != nil {
		return false, fmt.Errorf("error checking printer count: %v", err)
	}
	if printerCount >= 28 {
		return false, fmt.Errorf("maximum number of printers (28) already reached")
	}
	// --- END OF CHANGES ---

	// Check if printer ID already exists
	var existingId int
	err = database.DB.QueryRow("SELECT id FROM printers WHERE id = ?", request.Id).Scan(&existingId)
	if err == nil {
		// Row exists, printer ID is already taken
		return false, fmt.Errorf("printer with specified ID %d already exists", request.Id)
	} else if !errors.Is(err, sql.ErrNoRows) { // Use errors.Is for better error checking
		// Some other database error occurred
		return false, fmt.Errorf("error checking for existing printer ID: %v", err)
	}
	// No row exists with this ID, proceed with add

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
		} else {
			// Only commit if txErr is nil
			commitErr := tx.Commit()
			if commitErr != nil {
				// If commit fails, log it, but the primary error (txErr) might be more relevant if it exists
				log.Printf("Error committing transaction for AddPrinter: %v", commitErr)
				// Ensure txErr reflects the commit failure if it was previously nil
				if txErr == nil {
					txErr = commitErr
				}
			}
		}
	}()

	// Find the maximum rack_position for the given rack within the transaction
	var maxRackPosition sql.NullInt64
	// Use tx.QueryRow within the transaction
	err = tx.QueryRow("SELECT MAX(rack_position) FROM printers WHERE rack = ?", request.Rack).Scan(&maxRackPosition)
	// Check specifically for ErrNoRows, which is okay here (means first in rack)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		txErr = fmt.Errorf("error finding max rack position: %v", err)
		return false, txErr
	}

	// Calculate the next rack position
	var newRackPosition int
	if maxRackPosition.Valid {
		newRackPosition = int(maxRackPosition.Int64) + 1
	} else {
		newRackPosition = 1 // First printer in this rack
	}

	// Insert the new printer with the calculated rack_position
	insertSQL := `INSERT INTO printers (id, name, color, rack, rack_position, in_use, last_reserved_by, is_executive) values (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(insertSQL,
		request.Id,
		request.Name,
		request.Color,
		request.Rack,
		newRackPosition, // Use calculated position
		false,           // New printers are not in use
		nil,             // No one has reserved it yet
		request.Is_Executive)
	if err != nil {
		txErr = fmt.Errorf("error inserting new printer to DB: %v", err)
		return false, txErr
	}

	// If we reach here, txErr is nil, and the defer will commit.
	return true, nil // Return nil error on success
}

type UpdatePrinterRequest struct {
	Name         string `json:"name"`
	Color        string `json:"color"`
	Rack         int    `json:"rack"`
	RackPosition int    `json:"rack_position"` // Added rack position
	IsExecutive  bool   `json:"is_executive"`
}

// given printer id and attributes, update the printer.
// Checks if the target rack and position are already occupied by another printer.
func UpdatePrinter(id int, request UpdatePrinterRequest) (bool, error) {
	// Validate RackPosition
	if request.RackPosition <= 0 {
		return false, fmt.Errorf("invalid or missing rack_position: must be greater than 0")
	}
	// Validate Rack
	if request.Rack <= 0 {
		return false, fmt.Errorf("invalid or missing rack: must be greater than 0")
	}

	// Check if the printer to be updated exists
	var currentId int
	err := database.DB.QueryRow("SELECT id FROM printers WHERE id = ?", id).Scan(&currentId)
	if err == sql.ErrNoRows {
		return false, fmt.Errorf("printer with id %d does not exist", id)
	} else if err != nil {
		return false, fmt.Errorf("error checking if printer exists: %v", err)
	}

	// Check if the target rack and position is already occupied by *another* printer
	var conflictingId int
	err = database.DB.QueryRow("SELECT id FROM printers WHERE rack = ? AND rack_position = ? AND id != ?", request.Rack, request.RackPosition, id).Scan(&conflictingId)
	if err == nil {
		// A conflicting printer was found
		return false, fmt.Errorf("rack %d position %d is already occupied by printer ID %d", request.Rack, request.RackPosition, conflictingId)
	} else if err != sql.ErrNoRows {
		// An actual error occurred during the check
		return false, fmt.Errorf("error checking for conflicting printer position: %v", err)
	}
	// No conflict found, or the only printer at the target location is the one being updated (which is fine if rack/pos aren't changing)

	// Proceed with the update
	updateSQL := `UPDATE printers SET name = ?, color = ?, rack = ?, rack_position = ?, is_executive = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL,
		request.Name,
		request.Color,
		request.Rack,
		request.RackPosition, // Include rack position in update
		request.IsExecutive,
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
	// Include rack_position in the select query
	if err := database.DB.QueryRow("SELECT id, name, color, rack, rack_position, in_use, last_reserved_by, is_executive FROM printers WHERE id = ?", printerId).Scan(
		&printer.Id,
		&printer.Name,
		&printer.Color,
		&printer.Rack,
		&printer.Rack_Position, // Scan rack_position
		&printer.In_Use,
		&lastReservedBy,
		&printer.Is_Executive); err != nil {
		// Check if it's specifically a "no rows" error
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("printer with id %d not found", printerId)
		}
		return false, fmt.Errorf("failed to get printer details: %v", err)
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
		// This case should theoretically be caught by the initial printer check, but good to have defense in depth
		txErr = fmt.Errorf("no printer found with id: %d during update", printerId)
		return false, txErr
	}

	time_reserved := time.Now()
	time_complete := time.Now().Add(time.Duration(timeMins) * time.Minute)

	// Create reservation and add it to reservations table as an entry (within transaction)
	result, err = tx.Exec(
		"INSERT INTO reservations (printerid, userid, time_reserved, time_complete, is_active) values (?, ?, ?, ?, ?)",
		printerId,
		userId,
		time_reserved,
		time_complete,
		true)
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
		// Rollback happened in defer, but we still need to return the commit error
		return false, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Only turn on printer after successful transaction
	_, err = util.TurnOnPrinter(printerId)
	if err != nil {
		// Transaction was successful but printer failed to turn on
		// We should try to undo our changes
		undoErr := undoReservation(printerId, int(reservationId), userId, currentWeeklyMinutes)
		if undoErr != nil {
			log.Printf("CRITICAL: failed to undo reservation after printer turn on error: %v. Manual intervention may be required.", undoErr)
		} else {
			log.Printf("Reservation %d for printer %d successfully rolled back due to TurnOnPrinter failure.", reservationId, printerId)
		}
		return false, fmt.Errorf("error turning on printer: %v. Reservation has been rolled back", err)
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
	var txErr error
	defer func() {
		if txErr != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Set printer back to not in use
	_, err = tx.Exec("UPDATE printers SET in_use = FALSE WHERE id = ?", printerId)
	if err != nil {
		txErr = fmt.Errorf("failed to undo printer status: %v", err)
		return txErr
	}

	// Set reservation to inactive
	_, err = tx.Exec("UPDATE reservations SET is_active = FALSE WHERE id = ?", reservationId)
	if err != nil {
		txErr = fmt.Errorf("failed to undo reservation status: %v", err)
		return txErr
	}

	// Restore user's weekly minutes
	_, err = tx.Exec("UPDATE users SET weekly_minutes = ? WHERE id = ?", originalWeeklyMinutes, userId)
	if err != nil {
		txErr = fmt.Errorf("failed to restore user minutes: %v", err)
		return txErr
	}

	return nil // txErr is nil, commit will happen in defer
}

// turn off the printer, set the relevant printer as not in use, set the reservation to no longer be active
func CompleteReservation(printerId, reservationId int) {

	//Turn off the printer
	_, err := util.TurnOffPrinter(printerId)
	if err != nil {
		log.Printf("failed to turn off printer %d: %v", printerId, err)
		// Decide if we should proceed or retry later? For now, log and continue.
	}

	//Set as not in_use
	_, err = database.DB.Exec(
		"UPDATE printers SET in_use = FALSE WHERE id = ?",
		printerId,
	)
	if err != nil {
		log.Printf("failed to update printer %d status to not in use: %v", printerId, err)
	}
	//Set the reservation as inactive
	_, err = database.DB.Exec(
		"UPDATE reservations SET is_active = FALSE WHERE id = ?",
		reservationId,
	)
	if err != nil {
		log.Printf("failed to update reservation %d status to inactive: %v", reservationId, err)
	}

	// Remove from the active manager map
	manager.Mutex.Lock()
	// Check if the reservation still exists in the map before deleting
	if res, ok := manager.Reservations[reservationId]; ok {
		// Optionally stop the timer if it hasn't fired yet (though it should have)
		res.Timer.Stop()
		delete(manager.Reservations, reservationId)
		log.Printf("Completed and removed reservation %d from active manager.", reservationId)
	} else {
		log.Printf("Reservation %d not found in active manager upon completion.", reservationId)
	}
	manager.Mutex.Unlock()
}

// Given a printerId, toggle its is_executive bool in the printers table
func SetPrinterExecutive(id int) error {

	var currentExecutiveness bool

	querySQL := `SELECT is_executive FROM printers WHERE id = ?`
	err := database.DB.QueryRow(querySQL, id).Scan(&currentExecutiveness)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("printer with id %d not found", id)
		}
		return fmt.Errorf("error getting printer executiveness from db: %v", err)
	}

	newExecutiveness := !currentExecutiveness

	updateSQL := `UPDATE printers SET is_executive = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL, newExecutiveness, id)
	if err != nil {
		return fmt.Errorf("error updating printer executiveness: %v", err)
	}

	log.Printf("Toggled is_executive for printer %d to %v", id, newExecutiveness)
	return nil
}

// GetPrintersByRackId returns all printers belonging to a specific rack, ordered by position.
func GetPrintersByRackId(rackId int) ([]models.Printer, error) {
	// Query printers for the given rackId, ordered by rack_position
	query := "SELECT id, name, color, rack, rack_position, in_use, last_reserved_by, is_executive FROM printers WHERE rack = ? ORDER BY rack_position ASC"

	rows, err := database.DB.Query(query, rackId)
	if err != nil {
		return nil, fmt.Errorf("query error fetching printers for rack %d: %v", rackId, err)
	}
	defer rows.Close()

	// Initialize as an empty, non-nil slice
	printers := []models.Printer{}
	for rows.Next() {
		var p models.Printer
		var lastReservedBy sql.NullString
		// Scan all fields including rack_position
		if err := rows.Scan(&p.Id, &p.Name, &p.Color, &p.Rack, &p.Rack_Position, &p.In_Use, &lastReservedBy, &p.Is_Executive); err != nil {
			// Return nil for the slice in case of a scan error, along with the error itself
			return nil, fmt.Errorf("scan error for rack %d: %v", rackId, err)
		}
		if lastReservedBy.Valid {
			p.Last_Reserved_By = lastReservedBy.String
		} else {
			p.Last_Reserved_By = "" // Ensure it's an empty string if NULL
		}
		printers = append(printers, p)
	}

	// Check for errors during row iteration
	if err = rows.Err(); err != nil {
		// Return nil for the slice in case of a rows error, along with the error itself
		return nil, fmt.Errorf("rows error for rack %d: %v", rackId, err)
	}

	// Return the (potentially empty) slice and a nil error
	return printers, nil
}

// DeletePrinter removes a printer by its ID after checking for active reservations.
// If no active reservations exist, it also deletes associated inactive reservation history.
func DeletePrinter(id int) (bool, error) {
	// 1. Check if the printer exists
	var currentId int
	err := database.DB.QueryRow("SELECT id FROM printers WHERE id = ?", id).Scan(&currentId)
	if errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("printer with id %d not found", id)
	} else if err != nil {
		return false, fmt.Errorf("error checking if printer exists: %v", err)
	}

	// 2. Check for ACTIVE reservations associated with this printer
	var activeReservationCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM reservations WHERE printerid = ? AND is_active = TRUE", id).Scan(&activeReservationCount)
	if err != nil {
		return false, fmt.Errorf("error checking for active reservations: %v", err)
	}

	if activeReservationCount > 0 {
		return false, fmt.Errorf("cannot delete printer %d: it has %d active reservation(s)", id, activeReservationCount)
	}

	// 3. Start Transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %v", err)
	}
	// Use defer with a named error to handle rollback/commit
	var txErr error
	defer func() {
		if txErr != nil {
			log.Printf("Rolling back transaction due to error: %v", txErr)
			tx.Rollback()
		} else {
			commitErr := tx.Commit()
			if commitErr != nil {
				log.Printf("Error committing transaction for DeletePrinter: %v", commitErr)
				// Ensure txErr reflects the commit failure if it was previously nil
				if txErr == nil {
					txErr = commitErr // This won't be returned directly but signals rollback failure if needed
				}
			}
		}
	}()

	// 4. Delete INACTIVE reservations associated with the printer (within transaction)
	// Since we already checked for active ones, all remaining reservations for this printer ID must be inactive.
	_, err = tx.Exec("DELETE FROM reservations WHERE printerid = ?", id)
	if err != nil {
		txErr = fmt.Errorf("error deleting reservation history for printer %d: %v", id, err)
		return false, txErr
	}
	log.Printf("Deleted reservation history for printer %d", id)

	// 5. Attempt to delete the printer (within transaction)
	result, err := tx.Exec("DELETE FROM printers WHERE id = ?", id)
	if err != nil {
		// Check if the error is a foreign key constraint violation (shouldn't happen now, but good practice)
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint && strings.Contains(sqliteErr.Error(), "FOREIGN KEY constraint failed") {
				txErr = fmt.Errorf("unexpected foreign key constraint when deleting printer %d after deleting reservations: %v", id, err)
				return false, txErr
			}
		}
		// Otherwise, it's some other database error
		txErr = fmt.Errorf("error deleting printer %d: %v", id, err)
		return false, txErr
	}

	// 6. Verify deletion
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Log this error but don't necessarily fail the operation if deletion likely succeeded
		log.Printf("Warning: could not verify rows affected after deleting printer %d: %v", id, err)
		// Don't set txErr here, as the delete likely worked.
	}

	if rowsAffected == 0 {
		// This shouldn't happen if the initial check passed, but good to double-check
		txErr = fmt.Errorf("failed to delete printer %d (rows affected: 0), it might have been deleted concurrently", id)
		return false, txErr
	}

	// If we reach here, txErr is nil, and the defer will commit.
	log.Printf("Successfully deleted printer %d and its reservation history", id)
	return true, nil
}
