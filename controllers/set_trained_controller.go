package controllers

import (
	"gin-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
