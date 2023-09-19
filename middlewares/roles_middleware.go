package middlewares

import (
	// "enconding/json"
	"fileserver/repository"
	"fmt"

	"github.com/gin-gonic/gin"
	// "os"
)

func RolesMiddleware(c *gin.Context) {
	fmt.Println("ok roles")
	username, exists := c.Get("username")
	fmt.Println(username)
	fmt.Println(exists)
	if exists {
		str := string(username.(string))
		user := repository.GetUser(str)
		fmt.Println(user)
		fmt.Println(user.Access)
		if user.Access == 999 || user.Access == 555 {
			c.Next()
			return
		}

	}
	c.Redirect(302, "/notallowed")
	return
}
