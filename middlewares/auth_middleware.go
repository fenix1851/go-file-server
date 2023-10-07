package middlewares

import (
	// "enconding/json"
	"fileserver/repository"
	"fmt"

	"github.com/gin-gonic/gin"
	// "os"
)

func AuthMiddleware(c *gin.Context) {
	// get access token from cookie
	access_token, err := c.Cookie("access_token")
	fmt.Println(access_token, "access token from cookie in auth middleware")
	if err != nil {
		fmt.Println(err, "no access token")
		_, err := c.Cookie("refresh_token")
		if err != nil {
			fmt.Println(err, "no refresh token")
			c.Redirect(302, "/login")
			return
		} else {
			fmt.Println(err, "redirect to refresh")
			c.Redirect(302, "/refresh")
		}
	}
	// get user from access token
	user, err := repository.GetUserByToken(access_token, "access")
	fmt.Println(user, "user from access token in auth middleware")
	if err != nil {
		c.Redirect(302, "/login")
		return
	}
	if user.Username == "" {
		c.Redirect(302, "/login")
		return
	}
	// set user in context
	c.Set("username", user.Username)
	c.Next()
}
