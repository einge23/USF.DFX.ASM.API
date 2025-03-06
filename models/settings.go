package models

import "time"

type Settings struct {
	TimeSettings TimeSettings `json:"time_settings"`
}

type TimeSettings struct {
	WeekdayPrintTime PrintTime `json:"weekday_print_time"`
	WeekendPrintTime PrintTime `json:"weekend_print_time"`
	DayStart time.Time `json:"day_start"`
	NightStart time.Time `json:"night_start"`
	DefaultUserWeeklyHours int `json:"default_user_weekly_hours"`

}

type PrintTime struct {
	DayMaxPrintHours int `json:"day_max_print_hours"`
	NightMaxPrintHours int `json:"night_max_print_hours"`
}