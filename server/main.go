package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"7linternational.com/rest-server/db"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		log.Println("Receiving request in /ping")

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/login", func(c *gin.Context) {
		log.Println("Login attempt")
		var request LoginRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			user, code, err := db.Login(request.Email, request.Password)

			if err != nil {
				c.JSON(int(code), gin.H{
					"error": err.Error(),
				})
			}

			c.JSON(http.StatusOK, gin.H{
				"id":        user.Id,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
			})
		}
	})

	r.POST("/user", func(c *gin.Context) {
		log.Println("Register attempt")
		var request RegisterRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			code, err := db.Register(request.FirstName, request.LastName, request.Email, request.Password)

			if err != nil {
				c.JSON(int(code), gin.H{
					"error": err.Error(),
				})
			}

			c.JSON(http.StatusOK, gin.H{})
		}

	})

	r.GET("/user/:id/recipes", func(c *gin.Context) {
		idParam := c.Param("id")

		id, err := strconv.ParseInt(idParam, 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		code, err, recipes := db.GetUserRecipes(id)

		if err != nil {
			c.JSON(int(code), gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
				"recipes": recipes,
			})
	})

	log.Println("Start listening")
	r.Run(":8080")
}