package models

import (
	"database/sql"
	"encoding/json"
)

type UserData struct {
	Id                   int          `json:"id"`
	Username             string       `json:"username"`
	Trained              bool         `json:"trained"`
	Admin                bool         `json:"admin"`
	Has_Executive_Access bool         `json:"has_executive_access"`
	Ban_Time_End         sql.NullTime `json:"-"`
	Weekly_Minutes       int          `json:"weekly_minutes"`
}

func (u UserData) MarshalJSON() ([]byte, error) {
    type Alias UserData
    
    // Create a temporary struct for marshaling
    aux := struct {
        Alias
        Ban_Time_End interface{} `json:"ban_time_end"`
    }{
        Alias: Alias(u),
    }
    
    // Set Ban_Time_End to null or the value based on Valid flag
    if u.Ban_Time_End.Valid {
        aux.Ban_Time_End = u.Ban_Time_End.Time
    } else {
        aux.Ban_Time_End = nil
    }
    
    return json.Marshal(aux)
}