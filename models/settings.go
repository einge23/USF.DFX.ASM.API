package models

type Settings struct {
	TimeSettings          TimeSettings `json:"time_settings"`
	UpToDate              bool         `json:"up_to_date"`
	MaxActiveReservations int          `json:"max_active_reservations"`
}

type TimeSettings struct {
	WeekdayPrintTime       PrintTime `json:"weekday_print_time"`
	WeekendPrintTime       PrintTime `json:"weekend_print_time"`
	DayStart               string    `json:"day_start"`
	NightStart             string    `json:"night_start"`
	DefaultUserWeeklyHours int       `json:"default_user_weekly_hours"`
}

type PrintTime struct {
	DayMaxPrintHours   int `json:"day_max_print_hours"`
	NightMaxPrintHours int `json:"night_max_print_hours"`
}
