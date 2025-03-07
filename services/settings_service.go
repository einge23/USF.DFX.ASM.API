package services

import (
	"fmt"
	"gin-api/database"
	"gin-api/models"
	"gin-api/util"
)

//get the settings from the global obj if its up to date, and if it isnt, update it
//by pulling the info out of the db
func GetSettings() (models.Settings, error) {

	var err error = nil          //no error by default
	if !util.Settings.UpToDate { //if not up to date, fetch settings from DB
		err = util.ImportSettingsFromDB()
	}
	return util.Settings, err //return error if it exists, still nil if no error
}

//same as models.Settings but no up to date bool
type SetSettingsRequest struct {
	TimeSettings models.TimeSettings `json:"time_settings"`
}

func SetSettings(request SetSettingsRequest) error {

	//update global obj
	util.Settings.TimeSettings.WeekdayPrintTime.DayMaxPrintHours = request.TimeSettings.WeekdayPrintTime.DayMaxPrintHours
	util.Settings.TimeSettings.WeekdayPrintTime.NightMaxPrintHours = request.TimeSettings.WeekdayPrintTime.NightMaxPrintHours
	util.Settings.TimeSettings.WeekendPrintTime.DayMaxPrintHours = request.TimeSettings.WeekendPrintTime.DayMaxPrintHours
	util.Settings.TimeSettings.WeekendPrintTime.NightMaxPrintHours = request.TimeSettings.WeekendPrintTime.NightMaxPrintHours
	util.Settings.TimeSettings.DayStart = request.TimeSettings.DayStart
	util.Settings.TimeSettings.NightStart = request.TimeSettings.NightStart
	util.Settings.TimeSettings.DefaultUserWeeklyHours = request.TimeSettings.DefaultUserWeeklyHours

	//if somehow util is not up to date yet, set it to true
	util.Settings.UpToDate = true

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
