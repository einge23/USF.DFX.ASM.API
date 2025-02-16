package controllers

import (
	"gin-api/services"

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