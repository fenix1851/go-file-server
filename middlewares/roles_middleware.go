package middlewares

import (
	// "enconding/json"
	"fileserver/repository"
	"fmt"

	"github.com/gin-gonic/gin"
	// "os"
)

func RoleMiddleware(c *gin.Context) {
	username, exists := c.Get("username")
	if exists {
		str := fmt.Sprintf("%d", username)
		user := repository.GetUser(str)
		if user.Access == 999 || user.Access == 555 {
			c.Next()
		}

	}
	c.Redirect(403, "/403")
	return
}
