package models

type UserData struct {
	Id                   int    `json:"id"`
	Username             string `json:"username"`
	Trained              bool   `json:"trained"`
	Admin                bool   `json:"admin"`
	Has_Executive_Access bool   `json:"has_executive_access"`
}