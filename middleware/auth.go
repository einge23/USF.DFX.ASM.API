package middleware

import (
	"gin-api/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        token, err := util.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

		claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        c.Set("userId", claims["userId"])
        c.Set("isAdmin", claims["isAdmin"])
        c.Set("isEgnLab", claims["isEgnLab"])
        c.Next()
    }
}

func AdminPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
        isAdmin, exists := c.Get("isAdmin")
        if !exists || !isAdmin.(bool) {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access Denied."})
            c.Abort()
            return
        }
        c.Next()
    }
}

func UserOwnershipPermission() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get the requesting user's ID from the JWT token
        requestingUserID := c.GetInt("userId")
        // Get the target user ID from URL parameter
        targetUserID, err := strconv.Atoi(c.Param("userID"))
        
        if err != nil {
            c.JSON(400, gin.H{"error": "Invalid user ID"})
            c.Abort()
            return
        }

        // Allow if user is admin or requesting their own data
        isAdmin := c.GetBool("isAdmin")
        if !isAdmin && requestingUserID != targetUserID {
            c.JSON(403, gin.H{"error": "Unauthorized access"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
