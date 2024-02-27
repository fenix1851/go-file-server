package handlers

import (
	"fileserver/repository"
	"fileserver/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "os"
	// "path/filepath"
)

func LoginHandler(c *gin.Context) {
	// check query method
	if c.Request.Method == "GET" {
		c.HTML(200, "login.html", gin.H{})
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
		if user.Username == "" {
			c.HTML(401, "login.html", gin.H{"Error": "User does not exist"})
			return
		}
		// check if password is correct
		hashedPassword := utils.HashPassword(password)

		if user.HashedPassword != hashedPassword {
			c.JSON(401, gin.H{"error": "Incorrect password"})
			return
		}
		// create token
		access_token, err1 := utils.CreateToken(user.Username, "access")
		refresh_token, err2 := utils.CreateToken(user.Username, "refresh")
		//you had unused err while creating access token
		if err1 != nil || err2 != nil {
			c.HTML(500, "login.html", gin.H{"Error": "Error creating token"})
			return
		}
		// update user
		user.AccessToken = access_token
		user.RefreshToken = refresh_token
		repository.UpdateUser(DB, user)
		if err != nil {
			c.HTML(500, "login.html", gin.H{"Error": "Error updating user"})
			return
		}
		// set cookie
		c.SetCookie("access_token", user.AccessToken, 60*60*24, "/", "", false, false)
		c.SetCookie("refresh_token", user.RefreshToken, 60*60*24*14, "/", "", false, false)
		c.SetCookie("username", user.Username, 60*60*24, "/", "", false, false)
		c.Redirect(302, "/")
		return
	}
}
