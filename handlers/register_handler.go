package handlers

import (
	"fileserver/repository"

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
		user := repository.GetUser(username)
		if user.Username != "" {
			c.HTML(401, "register.html", gin.H{"Error": "User already exists"})
			return
		}
		// create user
		user = repository.CreateUser(username, password)
		// set cookies
		c.SetCookie("access_token", user.AccessToken, 60*60*24, "/", "localhost", false, true)
		c.SetCookie("refresh_token", user.RefreshToken, 60*60*24*14, "/", "localhost", false, true)
		c.SetCookie("username", user.Username, 60*60*24, "/", "localhost", false, true)
		c.Redirect(302, "/login")
	}
}
