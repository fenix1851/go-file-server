package handlers

import (
	"fileserver/repository"
	// "fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "os"
	"strconv"
)

func AdminHandler(c *gin.Context) {
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

	users, err := repository.GetUsers(DB)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.HTML(200, "admin.html", gin.H{
		"Users": users,
	})
}

func PatchUser(c *gin.Context) {
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

	username := c.PostForm("username")
	access := c.PostForm("access")
	// convert access to int
	accessInt, err := strconv.Atoi(access)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	user, err := repository.GetUser(DB, username)
	if err != nil {
		c.JSON(500, gin.H{"Error": err})
	}
	if user.Access == 999 {
		c.Redirect(302, "/admin")
		return
	}
	user.Access = accessInt
	repository.UpdateUser(DB, user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(302, "/admin")
}
