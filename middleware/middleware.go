package middleware

import (
	"auth-system/interfaces"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware struct {
	repo interfaces.Repository
	secretKey string
	client *redis.Client
}

func NewAuthMiddleware(r interfaces.Repository, s string, c *redis.Client) (AuthMiddleware) {
	return AuthMiddleware{repo: r, secretKey: s, client: c}
}

func (a AuthMiddleware) Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID string
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")

		userID, err := a.client.Get(c.Request.Context(), tokenString).Result()

		if err == redis.Nil {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(a.secretKey), nil
			})

			if err != nil {
				c.JSON(401, gin.H{
					"error": "Unauthorized",
				})
				c.Abort()
				return
			}

			if !token.Valid {
				c.JSON(401, gin.H{
					"error": "Unauthorized",
				})
				c.Abort()
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims);

			if !ok {
				c.JSON(401, gin.H{
					"error": "Unauthorized",
				})
				c.Abort()
				return
			}

			tokenID := claims["id"].(float64)
			accessToken, err := a.repo.GetToken(int(tokenID))
			if err != nil {
				c.JSON(401, gin.H{
					"error": "Unauthorized",
				})
				c.Abort()
				return
			}

			if accessToken.Revoked || accessToken.Expires_at.Before(time.Now()) {
				c.JSON(401, gin.H{
					"error": "Unauthorized",
				})
				c.Abort()
				return
			}

			userID = strconv.Itoa(int(accessToken.UserID))
			fmt.Println("Did not use redis")
		}
		
		c.Set("UserID", userID)
		c.Next()
	}
}