package controllers

import (
	"gin-api/models"
	"gin-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handles the GetTimeSettings service. Binds JSON to expected format and returns any errors encountered.
func GetTimeSettings(c *gin.Context) {
	var settings models.TimeSettings
	settings, err := services.GetTimeSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// handles the SetTimeSettings service. Binds JSON to expected format and returns any errors encountered.
func SetTimeSettings(c *gin.Context) {
	var req services.SetSettingsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.SetTimeSettings(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, true)
}

// handles the GetPrinterSettings service. Binds JSON to expected format and returns any errors encountered.
func GetPrinterSettings(c *gin.Context) {
	printerSettings, err := services.GetPrinterSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, printerSettings)
}

// Request body for setting printer settings
type SetPrinterSettingsRequest struct {
	MaxActiveReservations int `json:"max_active_reservations"`
}

// handles the SetPrinterSettings service. Binds JSON to expected format and returns any errors encountered.
func SetPrinterSettings(c *gin.Context) {
	var req SetPrinterSettingsRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.SetPrinterSettings(req.MaxActiveReservations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, true)
}

// handles the ExportDbToUsb service. Binds JSON to expected format and returns any errors encountered.
func ExportDbToUsb(c *gin.Context) {
	var req services.ExportDbToUsbRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := services.ExportDbToUsb(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, true)
}
