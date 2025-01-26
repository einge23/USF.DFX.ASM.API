package controllers

import (
	"gin-api/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
    var request services.LoginRequest
    
    if err := c.BindJSON(&request); err != nil {
        log.Printf("Error Binding Json: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userData, err := services.Login(request)
    if err != nil {
        log.Printf("Error in Login Service: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, userData)
}