package handlers

import (
	"fileserver/database"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	// check query method
	if c.Request.Method == "GET" {
		c.HTML(200, "register.html", gin.H{})
	} else if c.Request.Method == "POST" {
		// parse form
		username := c.PostForm("username")
		password := c.PostForm("password")
		// check if user exists
		user, err := database.GetUser(database.DB, username)
		if err != nil{
			c.JSON(500, gin.H{"Error": err})
		}
		if user.Username != "" {
			c.HTML(401, "register.html", gin.H{"Error": "User already exists"})
			return
		}
		// create user
		user, err = database.CreateUser(database.DB, username, password, 0)
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
