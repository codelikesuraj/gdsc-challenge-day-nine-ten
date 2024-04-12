package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/codelikesuraj/gdsc-challenge-day-nine-ten/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func (uc *UserController) Register(c *gin.Context) {
	var count int64
	var user models.User

	if err := c.ShouldBind(&user); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "invalid input",
				"errors":  models.GetValidationErrs(ve),
			})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
			"errors":  nil,
		})
		return
	}

	// check if user exists
	result := uc.DB.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	switch {
	case result.Error != nil:
		log.Println(result.Error)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	case count > 0:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "user already exists"})
		return
	}

	// hash password
	if err := user.HashPassword(user.Password); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	// create user
	result = uc.DB.Create(&user)
	if result.Error != nil || result.RowsAffected < 1 {
		log.Println(result.Error)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (uc *UserController) Login(c *gin.Context) {
	var user models.User

	if err := c.ShouldBind(&user); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "invalid input",
				"errors":  models.GetValidationErrs(ve),
			})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
			"errors":  nil,
		})
		return
	}

	username, password := user.Username, user.Password

	// hash password
	if err := user.HashPassword(password); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	// check if user exists
	result := uc.DB.Where("username = ?", username).First(&user)
	switch {
	case result.Error == gorm.ErrRecordNotFound:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid credentials"})
		return
	case result.Error != nil:
		log.Println(result.Error)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	// check password
	err := user.CheckPassword(password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid credentials"})
		return
	}

	tokenPair, err := generateTokenPair(user)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tokenPair})
}

func (uc *UserController) RefreshToken(c *gin.Context) {
	refreshTokenInput := struct {
		RefreshToken string `binding:"required" json:"refresh_token"`
	}{}

	if err := c.ShouldBind(&refreshTokenInput); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "invalid input",
				"errors":  models.GetValidationErrs(ve),
			})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
			"errors":  nil,
		})
		return
	}

	tokenString := refreshTokenInput.RefreshToken
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

	var user models.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		err := uc.DB.First(&user, claims["sub"]).Error
		switch {
		case err == gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid user"})
			return
		case err != nil:
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}
	}

	tokenPair, err := generateTokenPair(user)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tokenPair})
}

func (uc *UserController) Validate(c *gin.Context) {
	auth, _ := c.Get("auth_id")
	c.JSON(http.StatusOK, gin.H{
		"message": "I am logged in!",
		"data":    auth,
	})
}

func generateTokenPair(user models.User) (map[string]string, error) {
	accessToken, err := generateToken(time.Minute*30, user)
	if err != nil {
		return map[string]string{}, err
	}

	refreshToken, err := generateToken(time.Hour, user)
	if err != nil {
		return map[string]string{}, err
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func generateToken(expiration time.Duration, user models.User) (string, error) {
	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(expiration).Unix(),
	})

	// sign and get the complete encoded token as a string using the secret key
	return token.SignedString([]byte("SECRET_KEY"))
}
