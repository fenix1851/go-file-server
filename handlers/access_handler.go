package handlers

import "github.com/gin-gonic/gin"

func AccessHandler(c *gin.Context) {
	c.HTML(403, "access.html", gin.H{"Error": "not enough access"})
	return
}
