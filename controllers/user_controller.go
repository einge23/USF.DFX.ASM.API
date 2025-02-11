package controllers

import (
	"gin-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var req services.CreateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	failed, err := services.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		return
	}

	if failed {
		c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully added"})
}
