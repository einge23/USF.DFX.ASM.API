package controllers

import (
	"fmt"
	"gin-api/models"
	"gin-api/services"
	"gin-api/util"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// handles the GetPrinters service. Binds JSON to expected format and returns any errors encountered.
func GetPrinters(c *gin.Context) {

	printers, err := services.GetPrinters()
	if err != nil {
		log.Printf("Error in GetPrinters Service: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, printers)
}

// GetPrintersByRackId handles fetching printers for a specific rack ID.
func GetPrintersByRackId(c *gin.Context) {
	rackId := util.GetInfoFromPath(c, "rackId")
	if rackId == -1 {
		// GetInfoFromPath already sends a response (likely 400 Bad Request)
		return
	}

	printers, err := services.GetPrintersByRackId(rackId)
	if err != nil {
		log.Printf("Error in GetPrintersByRackId Service for rack %d: %v", rackId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve printers for the specified rack"})
		return
	}

	// Return the list of printers (which might be empty if none found, that's okay)
	c.JSON(http.StatusOK, printers)
}

// handles the AddPrinter service. Binds JSON to expected format and returns any errors encountered.
func AddPrinter(c *gin.Context) {
	var req models.Printer
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()}) // More specific error
		return
	}

	success, err := services.AddPrinter(req)
	if err != nil {
		// Check for specific user-facing errors vs internal errors
		if err.Error() == "maximum number of printers (28) already reached" ||
			strings.Contains(err.Error(), "already exists") ||
			strings.Contains(err.Error(), "invalid printer ID") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			log.Printf("Error in AddPrinter Service: %v", err) // Log internal errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add printer due to an internal error"})
		}
		return
	}
	// 'success' boolean isn't strictly needed if error is nil, but we keep it for consistency
	if !success {
		// This case might be redundant if errors are handled properly above
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add printer for an unknown reason"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Printer %d added successfully", req.Id)}) // Use 201 Created
}

// handles the UpdatePrinter service. Binds JSON to expected format and returns any errors encountered.
// requires that the printerId is given at the end of the route.
func UpdatePrinter(c *gin.Context) {

	id := util.GetInfoFromPath(c, "printerID")
	if id == -1 {
		return
	}

	var req services.UpdatePrinterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	success, err := services.UpdatePrinter(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "status not found"})
		return
	}

	c.JSON(http.StatusOK, true)
}

// handles the ReservePrinter service. Binds JSON to expected format and returns any errors encountered.
func ReservePrinter(c *gin.Context) {
	var req services.ReservePrinterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	success, err := services.ReservePrinter(req.PrinterId, req.UserId, req.TimeMins)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Printer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Printer status updated"})
}

// handles the UpdatePrinter service. Binds JSON to expected format and returns any errors encountered.
// requires that the printerId is given at the end of the route.
func SetPrinterExecutive(c *gin.Context) {

	id := util.GetInfoFromPath(c, "printerID")
	if id == -1 {
		return
	}

	err := services.SetPrinterExecutive(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, true)
}

// handles the DeletePrinter service.
// requires that the printerId is given at the end of the route.
func DeletePrinter(c *gin.Context) {
	id := util.GetInfoFromPath(c, "printerID")
	if id == -1 {
		// GetInfoFromPath already sends a response
		return
	}

	success, err := services.DeletePrinter(id)
	if err != nil {
		// Check for specific user-facing errors (not found, active reservations, history)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "active reservation") || strings.Contains(err.Error(), "past reservation history") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // 409 Conflict is appropriate here
		} else {
			log.Printf("Error in DeletePrinter Service for ID %d: %v", id, err) // Log internal errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete printer due to an internal error"})
		}
		return
	}

	// success bool might be redundant if error is nil
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete printer for an unknown reason"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Printer %d deleted successfully", id)})
}
