package services

import (
	"fmt"
	"gin-api/database"
	"gin-api/models"
	"gin-api/util"
)

// get the time settings from the global obj if its up to date.
// if it isnt up to date, update it by pulling the info out of the db
func GetTimeSettings() (models.TimeSettings, error) {

	var err error = nil          //no error by default
	if !util.Settings.TimeSettings.UpToDate { //if not up to date, fetch settings from DB
		err = util.ImportSettingsFromDB()
	}
	return util.Settings.TimeSettings, err //return error if it exists, still nil if no error
}

// same as models.Settings but no up to date bool
type SetSettingsRequest struct {
	TimeSettings models.TimeSettings `json:"time_settings"`
}

// directly set all time settings values both in the global obj and the database to avoid desync
func SetTimeSettings(request SetSettingsRequest) error {

	//update global obj
	util.Settings.TimeSettings.WeekdayPrintTime.DayMaxPrintHours = request.TimeSettings.WeekdayPrintTime.DayMaxPrintHours
	util.Settings.TimeSettings.WeekdayPrintTime.NightMaxPrintHours = request.TimeSettings.WeekdayPrintTime.NightMaxPrintHours
	util.Settings.TimeSettings.WeekendPrintTime.DayMaxPrintHours = request.TimeSettings.WeekendPrintTime.DayMaxPrintHours
	util.Settings.TimeSettings.WeekendPrintTime.NightMaxPrintHours = request.TimeSettings.WeekendPrintTime.NightMaxPrintHours
	util.Settings.TimeSettings.DayStart = request.TimeSettings.DayStart
	util.Settings.TimeSettings.NightStart = request.TimeSettings.NightStart
	util.Settings.TimeSettings.DefaultUserWeeklyHours = request.TimeSettings.DefaultUserWeeklyHours

	//if somehow UpToDate bool is not set yet, set it to true
	util.Settings.TimeSettings.UpToDate = true

	updateSQL := `UPDATE settings SET 
				day_max_print_hours_week = ?, night_max_print_hours_week = ?,
				day_max_print_hours_weekend = ?, night_max_print_hours_weekend = ?,
				day_start = ?, night_start = ?, default_user_weekly_hours = ?
				WHERE name = "default"`

	//update db
	_, err := database.DB.Exec(updateSQL,
		request.TimeSettings.WeekdayPrintTime.DayMaxPrintHours,
		request.TimeSettings.WeekdayPrintTime.NightMaxPrintHours,
		request.TimeSettings.WeekendPrintTime.DayMaxPrintHours,
		request.TimeSettings.WeekendPrintTime.NightMaxPrintHours,
		request.TimeSettings.DayStart,
		request.TimeSettings.NightStart,
		request.TimeSettings.DefaultUserWeeklyHours)
	if err != nil {
		return fmt.Errorf("error updating settings in db: %v", err)
	}
	return nil
}

//get the printer settings from global obj if it is up to date.
//If it is not up to date, import the settings from the DB and then get them.
func GetPrinterSettings() (models.PrinterSettings, error) {
	var err error = nil //no error by default
	if !util.Settings.PrinterSettings.UpToDate {
		err = util.ImportSettingsFromDB()
	}
	return util.Settings.PrinterSettings, err
}

//sets the printer settings passed in by the request. Currently the only printer
//setting is the max active reservations. Logic for other printer settings should be
//added here and request body should be added to.
func SetPrinterSettings(newMax int) error {
	if newMax <= 0 {
		return fmt.Errorf("max reservations must be a positive number")
	}
	//update global obj
	util.Settings.PrinterSettings.MaxActiveReservations = newMax

	//update in database
	updateSQL := `UPDATE settings SET max_active_reservations = ? WHERE name = "default"`
	database.DB.Exec(updateSQL, newMax)

	//raise upToDate flag for printerSettings
	util.Settings.PrinterSettings.UpToDate = true
	return nil
}
