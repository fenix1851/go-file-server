package handlers

import (
	"fileserver/repository"
	// "fmt"
	"github.com/gin-gonic/gin"
	// "os"
	"strconv"
)

func AdminHandler(c *gin.Context) {
	users, err := repository.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.HTML(200, "admin.html", gin.H{
		"Users": users,
	})
}

func PatchUser(c *gin.Context) {
	username := c.PostForm("username")
	access := c.PostForm("access")
	// convert access to int
	accessInt, err := strconv.Atoi(access)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	user := repository.GetUser(username)
	if user.Access == 999 {
		c.Redirect(302, "/admin")
		return
	}
	user.Access = accessInt
	repository.UpdateUser(user)
	c.Redirect(302, "/admin")
}
