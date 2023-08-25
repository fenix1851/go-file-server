package handlers

import (
	"fileserver/repository"
	"fileserver/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func RefreshHandler(c *gin.Context) {
	// get refresh token from cookie
	refresh_token, err := c.Cookie("refresh_token")
	if err != nil {
		c.Redirect(302, "/login")
		return
	}
	// get user from refresh token
	user, err := repository.GetUserByToken(refresh_token, "refresh")
	if err != nil {
		c.Redirect(302, "/login")
		return
	}
	if user.Username == "" || user.RefreshToken != refresh_token || user.RefreshToken == "" {
		c.Redirect(302, "/login")
		return
	}
	// generate new access token
	access_token, err := utils.CreateToken(user.Username, "access")
	if err != nil {
		fmt.Println(err)
	}
	// update user
	user.AccessToken = access_token
	repository.UpdateUser(user)
	// set cookies
	c.SetCookie("access_token", access_token, 60*60*24, "/", "localhost", false, true)
	c.SetCookie("refresh_token", user.RefreshToken, 60*60*24*14, "/", "localhost", false, true)
	c.SetCookie("username", user.Username, 60*60*24, "/", "localhost", false, true)
	c.Redirect(302, "/")
}
