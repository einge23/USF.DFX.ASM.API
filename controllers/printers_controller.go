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



func SetPrinterInUse(c *gin.Context) {
    var req services.SetInUseRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    success, err := services.SetPrinterInUse(req.PrinterId)
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