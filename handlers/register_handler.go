package handlers

import (
	"fileserver/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterHandler(c *gin.Context) {
	// check query method
	if c.Request.Method == "GET" {
		c.HTML(200, "register.html", gin.H{})
	} else if c.Request.Method == "POST" {
		dbInterface, exists := c.Get("db")
		if !exists {
			c.JSON(500, gin.H{"error": "DB not found in context"})
			return
		}

		DB, ok := dbInterface.(*gorm.DB)
		if !ok {
			c.JSON(500, gin.H{"error": "DB is not of type *gorm.DB"})
			return
		}
		// parse form
		username := c.PostForm("username")
		password := c.PostForm("password")
		// check if user exists
		user, err := repository.GetUser(DB, username)
		if err != nil {
			c.JSON(500, gin.H{"Error": err})
		}
		if user.Username != "" {
			c.HTML(401, "register.html", gin.H{"Error": "User already exists"})
			return
		}
		// create user
		user, err = repository.CreateUser(DB, username, password, 0)
		if err != nil {
			c.HTML(400, "register.html", gin.H{"Error": "Error while creating user"})
			return
		}
		// set cookies
		c.SetCookie("access_token", user.AccessToken, 60*60*24, "/", "localhost", false, true)
		c.SetCookie("refresh_token", user.RefreshToken, 60*60*24*14, "/", "localhost", false, true)
		c.SetCookie("username", user.Username, 60*60*24, "/", "localhost", false, true)
		c.Redirect(302, "/login")
	}
}
