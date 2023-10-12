package middlewares

import (
	// "enconding/json"
	"fileserver/repository"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "os"
)

func RolesMiddleware(minAccess int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if minAccess == 0 {
			c.Next()
			return
		}
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
		username, exists := c.Get("username")
		fmt.Println(username)
		fmt.Println(exists)
		if exists {
			str := string(username.(string))
			user, err := repository.GetUser(DB, str)
			if err != nil {
				c.JSON(500, gin.H{"Error": err})
			}
			fmt.Println(user)
			fmt.Println(user.Access)
			if user.Access > minAccess {
				c.Next()
				return
			}
		}
		c.Redirect(302, "/notallowed")
	}
}
