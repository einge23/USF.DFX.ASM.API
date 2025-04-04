package controllers

import (
	"gin-api/models"
	"gin-api/services"
	"gin-api/util"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

<<<<<<< HEAD
//handles the GetSettings service. Binds JSON to expected format and returns any errors encountered.
=======
// Get the settings by calling the service
>>>>>>> 5d0547a7c9d215e5f1ec0246cd23e176ff5b357b
func GetSettings(c *gin.Context) {
	var settings models.Settings
	settings, err := services.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}
<<<<<<< HEAD
//handles the SetSettings service. Binds JSON to expected format and returns any errors encountered.
=======

// Set the settings by calling the service and passing it the request body
>>>>>>> 5d0547a7c9d215e5f1ec0246cd23e176ff5b357b
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
	err = util.ExportTableToCSV("test.db", "reservations", outputPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exported to USB",
		"path":    outputPath,
	})
}
