package main

import (
	"fileserver/handlers"
	"fileserver/middlewares"
	"fileserver/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	gin := gin.Default()
	gin.LoadHTMLGlob("templates/*")
	gin.GET("/login", handlers.LoginHandler)
	gin.POST("/login", handlers.LoginHandler)
	gin.GET("/register", handlers.RegisterHandler)
	gin.POST("/register", handlers.RegisterHandler)
	gin.GET("/refresh", handlers.RefreshHandler)
	gin.Use(middlewares.AuthMiddleware)
	gin.GET("/", handlers.DirectoryHandler)
	gin.GET("/:path/*filepath", handlers.DirectoryHandler)
	utils.Init()
	gin.Run(":4001")
}
