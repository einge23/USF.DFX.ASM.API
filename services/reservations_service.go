package services

import (
	"fmt"
	"gin-api/database"
	"gin-api/models"
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
