package models

import (
	"sync"
	"time"
)

type Reservation struct {
	Id           int       `json:"id"`
	PrinterId    int       `json:"printer_id"`
	UserId       int       `json:"user_id"`
	Time_Reserved time.Time `json:"time_reserved"`
	Time_Complete time.Time `json:"time_complete"`
	Is_Active bool `json:"is_active"`
	Is_Egn_Reservation bool `json:"is_egn_reservation"`
	Timer *time.Timer
}

type ReservationManager struct {
	Reservations map[int]*Reservation
	Mutex sync.RWMutex
}

type ReservationDTO struct {
	Id                 int       `json:"id"`
	PrinterId          int       `json:"printer_id"`
	PrinterName        string    `json:"printer_name"` // Added
	UserId             int       `json:"user_id"`
	Username           string    `json:"username"` // Added
	Time_Reserved      time.Time `json:"time_reserved"`
	Time_Complete      time.Time `json:"time_complete"`
	Is_Active          bool      `json:"is_active"`
	Is_Egn_Reservation bool      `json:"is_egn_reservation"`
}