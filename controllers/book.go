package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/codelikesuraj/gdsc-challenge-day-nine-ten/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type BookController struct {
	DB *gorm.DB
}

func (b *BookController) GetAllBooks(c *gin.Context) {
	var books []models.Book

	a, _ := c.Get("auth")
	auth := a.(models.User)

	result := b.DB.Find(&books, "user_id = ?", auth.ID)
	if err := result.Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "oops, something went wrong - " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"books": books}})
}

func (b *BookController) GetABook(c *gin.Context) {
	var book models.Book

	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	a, _ := c.Get("auth")
	auth := a.(models.User)

	err := b.DB.First(&book, "id = ? AND user_id = ?", id, auth.ID).Error
	switch {
	case err == gorm.ErrRecordNotFound:
		c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	case err != nil:
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "oops, something went wrong - " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"book": book}})
}

func (b *BookController) CreateBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBind(&book); err != nil {
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

	a, _ := c.Get("auth")
	auth := a.(models.User)

	book.UserID = auth.ID

	result := b.DB.Create(&book)
	if result.Error != nil || result.RowsAffected < 1 {
		log.Println(result.Error)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "oops, something went wrong"})
	}

	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"book": book}})
}
