package main

import (
	"log"

	"github.com/codelikesuraj/gdsc-challenge-day-nine-ten/controllers"
	"github.com/codelikesuraj/gdsc-challenge-day-nine-ten/middlewares"
	"github.com/codelikesuraj/gdsc-challenge-day-nine-ten/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	PORT = "3000"
)

func main() {
	// initialize database
	var err error
	if DB, err = gorm.Open(sqlite.Open("bookstoreapi.db"), &gorm.Config{}); err != nil {
		log.Fatalln("Error connecting to database:", err.Error())
	}

	// run migrations
	if err := DB.AutoMigrate(&models.User{}, &models.Book{}); err != nil {
		log.Fatalln("Error running migrations:", err.Error())
	}

	r := gin.Default()

	BookController := controllers.BookController{DB: DB}
	UserController := controllers.UserController{DB: DB}
	AuthMiddleware := middlewares.Authenticated{DB: DB}

	r.POST("/register", UserController.Register)
	r.POST("/login", UserController.Login)
	r.POST("/refresh-token", UserController.RefreshToken)

	r.Group("", AuthMiddleware.Authenticate).
		GET("/books", BookController.GetAllBooks).
		GET("/books/:id", BookController.GetABook).
		POST("/books", BookController.CreateBook)

	log.Fatalln(r.Run(":" + PORT))
}
