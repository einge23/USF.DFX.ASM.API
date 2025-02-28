package models

import "database/sql"

type UserData struct {
	Id                   int          `json:"id"`
	Username             string       `json:"username"`
	Trained              bool         `json:"trained"`
	Admin                bool         `json:"admin"`
	Has_Executive_Access bool         `json:"has_executive_access"`
	Is_Egn_Lab           bool         `json:"is_egn_lab"`
	Ban_Time_End         sql.NullTime `json:"ban_time_end"`
	Weekly_Minutes       int          `json:"weekly_minutes"`
}