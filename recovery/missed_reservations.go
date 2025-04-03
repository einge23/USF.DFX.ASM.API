package recovery

import (
	"fmt"
	"gin-api/database"
	"gin-api/services"
	"gin-api/util"
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

	for _, r := range reservations {
		if r.timeComplete.Before(time.Now()) { //if reservation still is_active but its end time has passed, it was missed in downtime. Complete it.
			services.CompleteReservation(r.printerId, r.id)
		} else { //if end time has not passed yet, turn the printer back on.
			util.TurnOnPrinter(r.printerId)
		}
	}

	return true, nil
}
