package middleware

import (
	"gin-api/util"
	"net/http"
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
