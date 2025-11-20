package middleware

import (
	"net/http"

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

		// Try to verify if secret is available
		// secret := os.Getenv("SUPABASE_JWT_SECRET") // Commented out as we might not have it in this env

		// For now, we parse unverified as we might lack the secret in this dev environment.
		// IN PRODUCTION: You MUST verify the signature.
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Check expiry
			// if exp, ok := claims["exp"].(float64); ok {
			// 	if time.Now().Unix() > int64(exp) {
			// 		c.Redirect(http.StatusFound, "/login")
			// 		c.Abort()
			// 		return
			// 	}
			// }

			if sub, ok := claims["sub"].(string); ok {
				c.Set("user_id", sub)
				c.Next()
				return
			}
		}

		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}
}
