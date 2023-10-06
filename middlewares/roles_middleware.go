package middlewares

import (
	// "enconding/json"
	"fileserver/database"
	"fmt"

	"github.com/gin-gonic/gin"
	// "os"
)

func RolesMiddleware(minAccess int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if minAccess == 0 {
			c.Next()
			return
		}
		username, exists := c.Get("username")
		fmt.Println(username)
		fmt.Println(exists)
		if exists {
			str := string(username.(string))
			user, err := database.GetUser(database.DB, str)
			if err != nil{
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
		return
	}
}
