package recovery

import (
	"fmt"
	"gin-api/database"
	"gin-api/services"
	"time"
)

func CompleteMissedReservations() (bool, error) {

	//pull id, printer_id, time_complete of all active reservations
	querySQL := `SELECT id, printerid, time_complete FROM reservations WHERE is_active = ?`
	rows, err := database.DB.Query(querySQL, true)
	if err != nil {
		return false, fmt.Errorf("failed to query failsafe reservations: %v", err)
	}

	var reservations []struct {
		id           int
		printerId    int
		timeComplete time.Time
	}

	var errRows error

	//iterate through all rows pulled out from query
	for rows.Next() {
		var r struct {
			id           int
			printerId    int
			timeComplete time.Time
		}
		errRows = rows.Scan(&r.id, &r.printerId, &r.timeComplete)
		if errRows != nil {
			return false, fmt.Errorf("error scanning row %d: %v", r.id, err)
		}
		reservations = append(reservations, r)
	}

	//if there was a rows error during the loop, return it
	if errRows = rows.Err(); errRows != nil {
		return false, fmt.Errorf("error occurred during rows iteration: %v", err)
	}
	rows.Close()

	//if a reservation is not marked as complete (is_active) but its end time has passed, it was missed in downtime. Complete it.
	for _, r := range reservations {
		if r.timeComplete.Before(time.Now()) {
			services.CompleteReservation(r.printerId, r.id)
		}
	}

	return true, nil
}
