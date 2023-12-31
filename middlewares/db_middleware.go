package middlewares

import (
	"fileserver/repository"

	"github.com/gin-gonic/gin"
)

func DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := repository.GetDBInstance()
		if err != nil {
			c.JSON(500, gin.H{"error": "couldnt create an db instance"})
			return
		}
		c.Set("db", db)
		c.Next()
	}
}
