package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codelikesuraj/gdsc-challenge-day-nine-ten/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Authenticated struct {
	DB *gorm.DB
}

func (a *Authenticated) Authenticate(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authorization cookie is not set"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("SECRET_KEY"), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user models.User

		err := a.DB.First(&user, claims["sub"]).Error
		switch {
		case err == gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid user"})
			return
		case err != nil:
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		c.Set("auth", user)

		c.Next()

		return
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "bad"})
}
