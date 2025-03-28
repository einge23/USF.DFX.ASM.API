package controllers

import (
	"gin-api/models"
	"gin-api/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get the settings by calling the service
func GetSettings(c *gin.Context) {
	var settings models.Settings
	settings, err := services.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// Set the settings by calling the service and passing it the request body
func SetSettings(c *gin.Context) {

	log.Println(">>> SetSettings called")

	var req services.SetSettingsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(">>> JSON succesfully parsed", req)

	err := services.SetSettings(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		return
	}

	log.Println(">>> SetSettings completed succesfully")
	c.JSON(http.StatusOK, true)
}
