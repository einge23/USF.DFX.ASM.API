package controllers

import (
	"gin-api/services"
	"gin-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

//handles the CreateUser service. Binds JSON to expected format and returns any errors encountered.
func CreateUser(c *gin.Context) {
	var req services.CreateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := services.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
		return
	}

	c.JSON(http.StatusOK, true)
}

//handles the SetUserTrained service. Binds JSON to expected format and returns any errors encountered.
//requires that the userId is given at the end of the route.
func SetUserTrained(c *gin.Context) {

	id := util.GetInfoFromPath(c, "userID")
	if id == -1 {
		return
	}

	err := services.SetUserTrained(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, true)
}

//handles the GetUserReservations service. Binds JSON to expected format and returns any errors encountered.
//requires that the userId is given at the end of the route.
func GetUserReservations(c *gin.Context) {
	id := util.GetInfoFromPath(c, "userID")
	if id == -1 {
		return
	}

	reservations, err := services.GetUserReservations(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

//handles the GetActiveUerReservations service. Binds JSON to expected format and returns any errors encountered.
//requires that the userId is given at the end of the route.
func GetActiveUserReservations(c *gin.Context) {
	id := util.GetInfoFromPath(c, "userID")
	if id == -1 {
		return
	}

	reservations, err := services.GetActiveUserReservations(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(reservations) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No active reservations found"})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

type GetUserRequest struct {
	ScannerMessage string `json:"scanner_message" binding:"required"`
}

//handles the GetUserById service. Binds JSON to expected format and returns any errors encountered.
func GetUserById(c *gin.Context) {
	var req GetUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Scanner message is required"})
		return
	}

	parsedData, err := util.ParseScannerString(req.ScannerMessage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scanner message"})
		return
	}

	userData, err := services.GetUserById(parsedData.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userData)
}

//handles the SetUserExecutiveAccess service. Binds JSON to expected format and returns any errors encountered.
//requires that the userId is given at the end of the route.
func SetUserExecutiveAccess(c *gin.Context) {

	id := util.GetInfoFromPath(c, "userID")
	if id == -1 {
		return
	}

	err := services.SetUserExecutiveAccess(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error:": err.Error()})
		return
	}

	c.JSON(http.StatusOK, true)
}

//handles the AddUserWeeklyMinutes service. Binds JSON to expected format and returns any errors encountered.
//requires that the userId is given at the end of the route.
func AddUserWeeklyMinutes(c *gin.Context) {

	id := util.GetInfoFromPath(c, "userID")
	if id == -1 {
		return
	}

	var req services.AddUserWeeklyMinutesRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.AddUserWeeklyMinutes(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error:": err.Error()})
		return
	}

	c.JSON(http.StatusOK, true)
}

//handles the SetUserBanTime service. Binds JSON to expected format and returns any errors encountered.
//requires that the userId is given at the end of the route.
func SetUserBanTime(c *gin.Context) {

	id := util.GetInfoFromPath(c, "userID")
	if id == -1 {
		return
	}

	var req services.SetUserBanTimeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.SetUserBanTime(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error:": err.Error()})
		return
	}

	c.JSON(http.StatusOK, true)
}

//handles the GetUserWeeklyMinutes service. Binds JSON to expected format and returns any errors encountered.
//requires that the userId is given at the end of the route.
func GetUserWeeklyMinutes(c *gin.Context) {
	id := util.GetInfoFromPath(c, "userID")
	if id == -1 {
		return
	}

	minutes, err := services.GetUserWeeklyMinutes(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, minutes)
}
