package controllers

import (
	"gin-api/services"
	"gin-api/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var req services.CreateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	failed, err := services.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		return
	}

	if failed {
		c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully added"})
}

func SetUserTrained(c *gin.Context) {
	//get userID from path
	userID := c.Param("userID")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	//create request with id used in the specified path
	req := services.SetUserTrainedRequest{UserToTrain: id}
	failed, err := services.SetUserTrained(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	if failed {
		c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User training status successfully updated"})
}

func GetUserReservations(c *gin.Context) {
	userId := c.Param("userID")
	id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
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
	userId := c.Param("userID")
	id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
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