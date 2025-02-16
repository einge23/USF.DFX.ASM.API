package controllers

import (
	"gin-api/services"
	"gin-api/util"
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

    userData, tokenPair, err := services.Login(request)
    if err != nil {
        log.Printf("Error in Login Service: %v", err)
        switch err {
            case services.ErrorUserNotFound:
                c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            case services.ErrorNotTrained:
                c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    
    c.SetCookie("refresh_token", tokenPair.RefreshToken, 60*60*24*7, "/", "", false, true)

    c.JSON(http.StatusOK, gin.H{
        "user": userData,
        "access_token": tokenPair.AccessToken,
    })
}

func RefreshToken(c *gin.Context) {
    refreshToken := c.GetHeader("Refresh-Token")
    if refreshToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
        return
    }

    newAccessToken, err := util.RefreshAccessToken(refreshToken)
    if err != nil {
        switch err {
        case util.ErrExpiredToken:
            c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
        case util.ErrInvalidToken:
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid refresh token"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to refresh token"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token": newAccessToken,
    })
}