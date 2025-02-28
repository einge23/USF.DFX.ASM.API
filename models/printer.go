package models

type Printer struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Color            string `json:"color"`
	Rack             int    `json:"rack"`
	In_Use           bool   `json:"in_use"`
	Last_Reserved_By string `json:"last_reserved_by"`
	Is_Executive     bool   `json:"is_executive"`
	Is_Egn_Printer   bool   `json:"is_egn_printer"`
}