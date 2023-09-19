package middlewares

import (
	// "enconding/json"
	"fileserver/repository"
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
			user := repository.GetUser(str)
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
