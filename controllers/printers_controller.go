package controllers

import (
	"gin-api/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPrinters(c *gin.Context) {
	printers, err := services.GetPrinters()
	if err != nil {
		log.Printf("Error in GetPrinters Service: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, printers)
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
	var req services.SetPrinterExecRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	failed, err := services.SetPrinterExecutive(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		return
	}

	if failed {
		c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "printer executive status successfully changed"})
}
