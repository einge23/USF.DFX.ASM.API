package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetInfoFromPath(c *gin.Context, requestedInfo string) int {
	infoString := c.Param(requestedInfo)
	info, err := strconv.Atoi(infoString)
	if err != nil {
		errorMSG := fmt.Sprintf("Invalid %s", requestedInfo)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMSG})
		return -1
	}
	return info
}