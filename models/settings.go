package models

//master settings struct
type Settings struct {
	TimeSettings		TimeSettings `json:"time_settings"`
	PrinterSettings		PrinterSettings `json:"printer_settings"`
}

//time settings struct
type TimeSettings struct {
	WeekdayPrintTime       PrintTime `json:"weekday_print_time"`
	WeekendPrintTime       PrintTime `json:"weekend_print_time"`
	DayStart               string    `json:"day_start"`
	NightStart             string    `json:"night_start"`
	DefaultUserWeeklyHours int       `json:"default_user_weekly_hours"`
	UpToDate               bool      `json:"up_to_date"`
}

//substruct for time settings
type PrintTime struct {
	DayMaxPrintHours   int `json:"day_max_print_hours"`
	NightMaxPrintHours int `json:"night_max_print_hours"`
}

//printer settings struct
type PrinterSettings struct {
	MaxActiveReservations int 	`json:"max_active_reservations"`
	UpToDate              bool	`json:"up_to_date"`
}
