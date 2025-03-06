package controllers

import (
	"gin-api/models"
	"gin-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	var settings models.Settings
	settings, err := services.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}