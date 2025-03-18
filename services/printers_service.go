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

//return all printers by rack, as serialized JSON
func GetPrinters() ([]models.Printer, error) {
	rows, err := database.DB.Query("SELECT id, name, color, rack, in_use, last_reserved_by, is_executive FROM printers order by rack asc")
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

//Add printer to the db by ensuring that a printer of that ID doesn't exist. ID correlates to physical plug.
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
			return false, fmt.Errorf("error inserting new printer to DB")
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
	Name 			string	`json:"name"`
	Color 			string	`json:"color"`
	Rack 			int		`json:"rack"`
	IsExecutive 	bool	`json:"is_executive"`
	IsEgnPrinter 	bool	`json:"is_egn_printer"`
}

//change name, color, rack, is_executive, is_egn, of an existing printer
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
		return false, fmt.Errorf("error updating printer in DB")
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

	//Turn on printer
	_, err := util.TurnOnPrinter(printerId)
	if err != nil {
		return false, fmt.Errorf("error turning on printer: %v", err)
	}

	//Set printer as 'in use' in the database
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

	time_reserved := time.Now()
	time_complete := time.Now().Add(time.Duration(timeMins) * time.Minute)

	//create reservation and add it to reservations table as an entry
	result, err = database.DB.Exec(
		"INSERT INTO reservations (printerid, userid, time_reserved, time_complete, is_active, is_egn_reservation) values (?, ?, ?, ?, ?, ?)",
		printerId,
		userId,
		time_reserved,
		time_complete,
		true,
		printer.Is_Egn_Printer)
	if err != nil {
		return false, fmt.Errorf("failed to insert reservation: %v", err)
	}

	reservationId, err := result.LastInsertId()
	if err != nil {
		return false, fmt.Errorf("failed to get reservation id: %v", err)
	}

	//get user's weeklyMinutes
	var currentWeeklyMinutes int
	querySQL := `SELECT weekly_minutes FROM users WHERE id = ?`
	err = database.DB.QueryRow(querySQL, userId).Scan(&currentWeeklyMinutes)
	if err != nil {
		return false, fmt.Errorf("error getting user weekly minutes from db: %v", err)
	}

	//subtract reservation's duration from weeklyMinutes
	newWeeklyMinutes := currentWeeklyMinutes - timeMins
	updateSQL := `UPDATE users SET weekly_minutes = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL, newWeeklyMinutes, userId)
	if err != nil {
		return false, fmt.Errorf("error subtracting minutes from user")
	}

	//set up timer to complete/end the reservation. calls 'completeReservation()' function when the timer is up.
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
		completeReservation(printerId, int(reservationId))
	}()

	return true, nil
}

// turn off the printer, set the relevant printer as not in use, set the reservation to no longer be active
func completeReservation(printerId, reservationId int) {

	//Turn off the printer
	_, err := util.TurnOffPrinter(printerId)
	if err != nil {
		log.Printf("failed to turn off printer: %v", err)
	}

	_, err = database.DB.Exec(
		"UPDATE printers SET in_use = FALSE WHERE id = ?",
		printerId,
	)
	if err != nil {
		log.Printf("failed to update printer: %v", err)
	}

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
