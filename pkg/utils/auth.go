package utils

import (
	"Backend/internal/services"
	"context"
	"errors"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
)

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDRaw, _ := c.Get("userID")
	userID := userIDRaw.(uuid.UUID)
	if userID == uuid.Nil {
		return uuid.Nil, errors.New("user id is nil")
	}

	return userID, nil
}

func GetUserIDFromToken(tokenString, secretKey string) (uuid.UUID, error) {
	claims, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return uuid.Nil, err
	}

	return claims.UserID, nil
}

func ExtractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}

	return "", errors.New("invalid authorization header")
}

func ValidateToken(tokenString, secretKey string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, &CustomError{
			ErrorResponse: ErrorResponse{
				Errors: []ErrorDetail{
					{
						Status:  http.StatusBadRequest,
						Message: "The provided token is invalid",
					},
				},
			},
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		storedToken, err := RetrieveTokenFromRedis(claims.UserID)
		if err != nil {
			return nil, &CustomError{
				ErrorResponse: ErrorResponse{
					Errors: []ErrorDetail{
						{
							Status:  http.StatusInternalServerError,
							Message: "Cannot retrieve token from redis",
						},
					},
				},
			}
		}

		if tokenString != storedToken {
			return nil, &CustomError{
				ErrorResponse: ErrorResponse{
					Errors: []ErrorDetail{
						{
							Status:  http.StatusUnauthorized,
							Message: "Token mismatch",
						},
					},
				},
			}
		}

		IsRevoked, err := IsTokenRevoked(tokenString)
		if err != nil {
			return nil, &CustomError{
				ErrorResponse: ErrorResponse{
					Errors: []ErrorDetail{
						{
							Status:  http.StatusInternalServerError,
							Message: "Cannot check if token is revoked",
						},
					},
				},
			}
		}

		if IsRevoked {
			return nil, &CustomError{
				ErrorResponse: ErrorResponse{
					Errors: []ErrorDetail{
						{
							Status:  http.StatusUnauthorized,
							Message: "Token is revoked",
						},
					},
				},
			}
		}

		return claims, nil
	}

	return nil, &CustomError{
		ErrorResponse: ErrorResponse{
			Errors: []ErrorDetail{
				{
					Status:  http.StatusUnauthorized,
					Message: "Unauthorized",
				},
			},
		},
	}
}

func ExtractUserIDAndCheckPermission(c *gin.Context, permissionType string) (uuid.UUID, error) {
	token, err := ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return uuid.UUID{}, err
	}

	userID, err := GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return uuid.UUID{}, err
	}

	hasPermission, err := services.PermissionService.CheckPermission(services.PermissionService{}, context.Background(), userID, permissionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return uuid.UUID{}, err
	} else if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": []string{"You don't have permission to perform this action"}})
		return uuid.UUID{}, errors.New("permission denied")
	}

	return userID, nil
}
