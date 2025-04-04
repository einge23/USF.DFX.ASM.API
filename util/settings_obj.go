package util

import (
	"fmt"
	"gin-api/database"
	"gin-api/models"
)

// Global variable referenced by other packages to get/set the settings. Should be
// faster than fetching from db every time
var Settings models.Settings

// Gets the settings from the database and stores them in the 'Settings' global obj
func ImportSettingsFromDB() error {
	querySQL := `SELECT day_max_print_hours_week, night_max_print_hours_week,
						day_max_print_hours_weekend, night_max_print_hours_weekend,
						day_start, night_start, default_user_weekly_hours
						FROM settings WHERE name = "default"`
	err := database.DB.QueryRow(querySQL).Scan(
		&Settings.TimeSettings.WeekdayPrintTime.DayMaxPrintHours,
		&Settings.TimeSettings.WeekdayPrintTime.NightMaxPrintHours,
		&Settings.TimeSettings.WeekendPrintTime.DayMaxPrintHours,
		&Settings.TimeSettings.WeekendPrintTime.NightMaxPrintHours,
		&Settings.TimeSettings.DayStart,
		&Settings.TimeSettings.NightStart,
		&Settings.TimeSettings.DefaultUserWeeklyHours)
	if err != nil {
		return fmt.Errorf("error getting settings from db: %v", err)
	}

	Settings.UpToDate = true
	return nil
}
