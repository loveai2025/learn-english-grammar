package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		secret := os.Getenv("SUPABASE_JWT_SECRET")
		var token *jwt.Token

		if secret != "" {
			// Verify if secret is available
			token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})
		} else {
			// Fallback to unverified if secret not set (DEV ONLY)
			// IN PRODUCTION: You MUST set SUPABASE_JWT_SECRET and verify the signature.
			token, _, err = new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		}

		if err != nil {
			respondUnauthorized(c)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok {
				c.Set("user_id", sub)
				c.Next()
				return
			}
		}

		respondUnauthorized(c)
	}
}

func respondUnauthorized(c *gin.Context) {
	// Check if it's an API request
	if c.Request.Header.Get("Accept") == "application/json" || len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	} else {
		c.Redirect(http.StatusFound, "/login")
	}
	c.Abort()
}
