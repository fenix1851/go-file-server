package handlers

import (
	// "os"
	"github.com/gin-gonic/gin"
)

func TestHandler(c *gin.Context) {
	c.HTML(200, "test.html", gin.H{"title": "Test Page"})
}
