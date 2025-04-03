package controllers

import (
	"gin-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetActiveReservations(c *gin.Context) {
	reservations, err := services.GetActiveReservations()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if reservations == nil {
		c.JSON(200, []interface{}{})
		return
	}

	c.JSON(200, reservations)
}

func CancelActiveReservation(c *gin.Context) {
	var req services.CancelActiveReservationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Internal server error:": err.Error()})
		return
	}

	_, err := services.CancelActiveReservation(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error:": err.Error()})
		return
	}

	c.JSON(http.StatusOK, true)
}
