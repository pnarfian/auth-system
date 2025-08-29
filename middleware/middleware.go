package middleware

import (
	"auth-system/interfaces"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthMiddleware struct {
	repo interfaces.Repository
	secretKey string
}

func NewAuthMiddleware(r interfaces.Repository, s string) (AuthMiddleware) {
	return AuthMiddleware{repo: r, secretKey: s}
}

func (a AuthMiddleware) Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("SecretKey: " + a.secretKey)
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
				"test": "no bearer",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
      return []byte(a.secretKey), nil
		})

		if err != nil {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
				"asda": "xcvxv",
				"err": err.Error(),
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
				"asdsa": "qweqwe",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims);

		if !ok {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
				"broo": "asdad",
			})
			c.Abort()
			return
		}

		tokenID := claims["id"].(float64)
		accessToken, err := a.repo.GetToken(int(tokenID))
		
		if err != nil {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
				"test": "mamamam",
			})
			c.Abort()
			return
		}

		if accessToken.Revoked || accessToken.Expires_at.Before(time.Now()) {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
				"test": "blabliblu",
			})
			c.Abort()
			return
		}
		
		c.Set("tokenID", tokenID)
		c.Set("UserID", accessToken.UserID)
		c.Next()
	}
}