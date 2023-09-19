package main

import (
	"fileserver/handlers"
	"fileserver/middlewares"
	"fileserver/startup"

	"github.com/gin-gonic/gin"
)

func main() {
	gin := gin.Default()
	gin.LoadHTMLGlob("templates/*")
	gin.GET("/login", handlers.LoginHandler)
	gin.POST("/login", handlers.LoginHandler)
	gin.GET("/access", handlers.AccessHandler)
	gin.GET("/register", handlers.RegisterHandler)
	gin.POST("/register", handlers.RegisterHandler)
	gin.GET("/refresh", handlers.RefreshHandler)
	gin.GET("/notallowed", handlers.AccessHandler)
	gin.Use(middlewares.AuthMiddleware)
	gin.Use(middlewares.RolesMiddleware)
	gin.GET("/", handlers.DirectoryHandler)
	gin.GET("/:path/*filepath", handlers.DirectoryHandler)
	startup.Init()
	gin.Run(":4001")
}
