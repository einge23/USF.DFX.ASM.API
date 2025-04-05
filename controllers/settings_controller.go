package controllers

import (
	"gin-api/models"
	"gin-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handles the GetSettings service. Binds JSON to expected format and returns any errors encountered.
func GetSettings(c *gin.Context) {
	var settings models.Settings
	settings, err := services.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// handles the SetSettings service. Binds JSON to expected format and returns any errors encountered.
func SetSettings(c *gin.Context) {
	var req services.SetSettingsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.SetSettings(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, true)
}
