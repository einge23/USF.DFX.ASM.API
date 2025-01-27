package services

import (
	"fmt"
	"gin-api/database"
	"gin-api/models"
	"time"
)

func GetPrinters() ([]models.Printer, error) {
	rows, err := database.DB.Query("SELECT id, name, color, rack, in_use FROM printers")
    if err != nil {
        return nil, fmt.Errorf("query error: %v", err)
    }
    defer rows.Close()

	var printers []models.Printer
	for rows.Next() {
		var p models.Printer
        if err := rows.Scan(&p.Id, &p.Name, &p.Color, &p.Rack, &p.In_Use); err != nil {
            return nil, fmt.Errorf("scan error: %v", err)
		}
		printers = append(printers, p)
	}

	if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %v", err)
	}
	return printers, nil
}

type ReservePrinterRequest struct {
    PrinterId int `json:"printer_id"`
	UserId int `json:"user_id"`
	TimeMins int `json:"time_mins"`
}

func ReservePrinter(printerId int, userId int, timeMins int) (bool, error) {
	var user models.UserData
	if err := database.DB.QueryRow("SELECT username FROM users WHERE id = ?", userId).Scan(
		&user.Username); err != nil {
		return false, fmt.Errorf("failed to get username: %v", err)
	}

	var printer models.Printer
	if err := database.DB.QueryRow("SELECT id, name, color, rack, in_use FROM printers WHERE id = ?", printerId).Scan(
		&printer.Id,
		&printer.Name,
		&printer.Color,
		&printer.Rack,
		&printer.In_Use); err != nil {
		return false, fmt.Errorf("failed to get printer: %v", err)
	}

	if printer.In_Use {
		return false, fmt.Errorf("printer is already in use")
	}
 	result, err := database.DB.Exec(
		"UPDATE printers SET in_use = TRUE, last_reserved_by = ? WHERE id = ?",
		user.Username,
		printerId,
	)
	if err != nil {
		return false, fmt.Errorf("failed to update printer: %v", err)
	}

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return false, fmt.Errorf("failed to get affected rows: %v", err)
    }

    if rowsAffected == 0 {
        return false, fmt.Errorf("no printer found with id: %d", printerId)
    }

	go func() {
		time.Sleep(time.Duration(timeMins) * time.Minute)
		_, err := database.DB.Exec(
            "UPDATE printers SET in_use = FALSE WHERE id = ?",
            printerId,
        )
		if err != nil {
			fmt.Printf("failed to release printer: %v", err)
		}
	}()

    return true, nil
}