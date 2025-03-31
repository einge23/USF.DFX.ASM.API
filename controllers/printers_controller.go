package controllers

import (
	"gin-api/models"
	"gin-api/services"
	"gin-api/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPrinters(c *gin.Context) {
	_, exists := c.Get("isEgnLab")
	if !exists {
		log.Printf("Error: No claims found in request context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	isEgnLab := c.GetBool("isEgnLab")

	printers, err := services.GetPrinters(isEgnLab)
	if err != nil {
		log.Printf("Error in GetPrinters Service: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, printers)
}

func AddPrinter(c *gin.Context) {
	var req models.Printer
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	success, err := services.AddPrinter(req)
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
