package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetInfoFromPath(c *gin.Context, requestedInfo string) int {
	infoString := c.Param(requestedInfo)
	info, err := strconv.Atoi(infoString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return -1
	}
	return info
}