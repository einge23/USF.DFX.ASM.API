package models

type Printer struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	Rack   int    `json:"rack"`
	In_Use bool   `json:"in_use"`
}