package controllers

import (
	"gin-api/services"
	"gin-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func SetUserTrained(c *gin.Context) {

	id := GetInfoFromPath(c, "userID")
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

func GetUserReservations(c *gin.Context) {
	id := GetInfoFromPath(c, "userID")
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

func GetActiveUserReservations(c *gin.Context) {
	id := GetInfoFromPath(c, "userID")
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

func SetUserExecutiveAccess(c *gin.Context) {

	id := GetInfoFromPath(c, "userID")
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
