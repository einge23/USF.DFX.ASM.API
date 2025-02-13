package services

import (
	"fmt"
	"gin-api/database"
)

type SetUserTrainedRequest struct {
	UserToTrain int
}

func SetUserTrained(setUserTrainedRequest SetUserTrainedRequest) (bool, error) {

	//get trained status for the user from db
	var trainedStatus bool
	querySQL := `SELECT has_training FROM users WHERE id = ?`
	err := database.DB.QueryRow(querySQL, setUserTrainedRequest.UserToTrain).Scan(&trainedStatus)
	if err != nil {
		return true, fmt.Errorf("error getting user from db: %v", err)
	}

	//toggle the trainedStatus
	newTrainedStatus := !trainedStatus

	//update the user's training status in the database
	updateSQL := `UPDATE users SET has_training = ? WHERE id = ?`
	_, err = database.DB.Exec(updateSQL, newTrainedStatus, setUserTrainedRequest.UserToTrain)
	if err != nil {
		return true, fmt.Errorf("error updating user training status: %v", err)
	}

	return false, nil
}
