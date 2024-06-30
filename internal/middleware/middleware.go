package middleware

import (
	"Backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"os"
)

func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := utils.ExtractTokenFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		claims, err := utils.ValidateToken(tokenString, os.Getenv("JWT_SECRET_KEY"))
		if err != nil {
			c.JSON(err.(*utils.CustomError).ErrorResponse.Errors[0].Status, gin.H{"success": false, "message": []string{err.Error()}})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}
