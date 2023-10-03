package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserIDFromContext(c *gin.Context) uuid.UUID {
	userID, _ := c.Get("userID")
	return userID.(uuid.UUID)
}
