package models

import (
	"sync"
	"time"
)

type Reservation struct {
	Id           int       `json:"id"`
	PrinterId    int       `json:"printer_id"`
	UserId       int       `json:"user_id"`
	TimeReserved time.Time `json:"time_reserved"`
	TimeComplete time.Time `json:"time_complete"`
	IsActive bool `json:"is_active"`
	Timer *time.Timer
}

type ReservationManager struct {
	Reservations map[int]*Reservation
	Mutex sync.RWMutex
}

type ReservationDTO struct {
	Id           int       `json:"id"`
	PrinterId    int       `json:"printer_id"`
	UserId       int       `json:"user_id"`
	TimeReserved time.Time `json:"time_reserved"`
	TimeComplete time.Time `json:"time_complete"`
	IsActive bool `json:"is_active"`
}