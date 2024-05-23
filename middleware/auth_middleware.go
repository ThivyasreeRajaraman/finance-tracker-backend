package middleware

import (
	"fmt"
	"net/http"
	"os"

	userhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/user"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == utils.EmptyString {
			utils.HandleError(c, http.StatusUnauthorized, "Authorization header missing", nil)
			c.Abort()
			return
		}

		token, err := parseToken(tokenString)
		if err != nil || !token.Valid {
			utils.HandleError(c, http.StatusUnauthorized, "Invalid token", err)
			c.Abort()
			return
		}
		if err := authenticateUserByToken(c, token); err != nil {
			utils.HandleError(c, http.StatusUnauthorized, "Invalid token claims", err)
			c.Abort()
			return
		}
	}
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
}

func authenticateUserByToken(c *gin.Context, token *jwt.Token) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}
	user, err := userhelper.FetchUserByClaims(claims)
	if err != nil {
		return fmt.Errorf("invalid user data")
	}
	c.Set("currentUser", user)
	c.Next()
	return nil
}
