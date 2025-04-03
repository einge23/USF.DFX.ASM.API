package services

import (
	"database/sql"
	"fmt"
	"gin-api/database"
	"gin-api/models"
	"time"
)

func GetActiveReservations() ([]models.ReservationDTO, error) {
	rows, err := database.DB.Query("SELECT id, printerId, time_reserved, time_complete, userId, is_active, is_egn_reservation FROM reservations WHERE is_active = 1")
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var reservations []models.ReservationDTO
	for rows.Next() {
		var r models.ReservationDTO
		if err := rows.Scan(&r.Id, &r.PrinterId, &r.Time_Reserved, &r.Time_Complete, &r.UserId, &r.Is_Active, &r.Is_Egn_Reservation); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		reservations = append(reservations, r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return reservations, nil
}

type CancelActiveReservationRequest struct {
	PrinterId     int `json:"printer_id"`
	ReservationId int `json:"reservation_id"`
}

// Cancel the reservation specified by the printerId and reservationId, refund the reservation's remaining time to the user
func CancelActiveReservation(request CancelActiveReservationRequest) (bool, error) {
	var userId int
	var isActive bool
	var timeComplete time.Time

	//pull userId and is_active from the reservation
	err := database.DB.QueryRow("SELECT userId, is_active, time_complete FROM reservations WHERE id = ?", request.ReservationId).Scan(&userId, &isActive, &timeComplete)

	if err == sql.ErrNoRows { //handle nonexistent reservation
		return false, fmt.Errorf("error cancelling reservation, no reservation of ID %d exists", request.ReservationId)
	} else if !isActive { //handle reservation that isn't active
		return false, fmt.Errorf("error cancelling reservation, the reservation requested for cancellation is not active")
	} else if err != nil { //handle all other errors from query
		return false, fmt.Errorf("error cancelling reservation: %v", err)
	}

	//get time that was left in the reservation
	timeToRefund := time.Until(timeComplete)
	if timeToRefund < 0 { //realistically shouldn't ever happen because we already checked if !isActive, more of a precaution than anything
		return false, fmt.Errorf("error cancelling reservation: the requested reservation is already over")
	}

	//convert time to minutes so its compatible with weekly_minutes db column
	minutesToRefund := int(timeToRefund.Minutes())

	var userWeeklyMinutes int
	//pull user's current weekly_minutes out of the database
	err = database.DB.QueryRow("SELECT weekly_minutes FROM users WHERE id = ?", userId).Scan(&userWeeklyMinutes)
	if err != nil {
		return false, fmt.Errorf("error getting weekly_minutes of the reservation's user: %v", err)
	}

	refundedUserWeeklyMinutes := userWeeklyMinutes + minutesToRefund //add
	updateSQL := `UPDATE users SET weekly_minutes = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL, refundedUserWeeklyMinutes, userId)
	if err != nil {
		return false, fmt.Errorf("error refunding weekly minutes to user: %v", err)
	}

	//now that we have refunded the cancellation without errors, remove the reservation formally
	CompleteReservation(request.PrinterId, request.ReservationId)
	return true, nil
}
