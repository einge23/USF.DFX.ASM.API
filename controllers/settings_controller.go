package controllers

import (
	"gin-api/models"
	"gin-api/services"
	"gin-api/util"
	"net/http"
	"path/filepath"

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

func ExportDBToUSB(c *gin.Context) {
	usbPath, err := util.FindUSBDrive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No USB drive found"})
		return
	}

	// Create full CSV path on the USB
	outputPath := filepath.Join(usbPath, "reservations_export.csv")

	// Export to that path
	err = util.ExportTableToCSV("test.db", "users", outputPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exported to USB",
		"path":    outputPath,
	})
}
