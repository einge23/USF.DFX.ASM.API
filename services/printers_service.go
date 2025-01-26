package services

import (
	"fmt"
	"gin-api/database"
	"gin-api/models"
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

type SetInUseRequest struct {
    PrinterId int `json:"printer_id"`
}

func SetPrinterInUse(printerId int) (bool, error) {
    result, err := database.DB.Exec("UPDATE printers SET in_use = TRUE WHERE id = ?", printerId)
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

    return true, nil
}