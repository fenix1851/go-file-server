package main

import (
	"fileserver/database"
	"fileserver/handlers"
	"fileserver/middlewares"
	"fileserver/startup"

	"github.com/gin-gonic/gin"
)

func main() {
	database.DBinitialize()
	gin := gin.Default()
	gin.LoadHTMLGlob("templates/*")
	gin.GET("/login", handlers.LoginHandler)
	gin.POST("/login", handlers.LoginHandler)
	gin.GET("/access", handlers.AccessHandler)
	gin.GET("/register", handlers.RegisterHandler)
	gin.POST("/register", handlers.RegisterHandler)
	gin.GET("/refresh", handlers.RefreshHandler)
	gin.GET("/notallowed", handlers.AccessHandler)
	gin.GET("/", middlewares.AuthMiddleware, middlewares.RolesMiddleware(1), handlers.DirectoryHandler)
	gin.POST("/", handlers.UploadHandler)
	gin.GET("/admin", middlewares.AuthMiddleware, middlewares.RolesMiddleware(998), handlers.AdminHandler)
	gin.POST("/admin/patch_user", middlewares.AuthMiddleware, middlewares.RolesMiddleware(998), handlers.PatchUser)
	gin.GET("/:path/*filepath", middlewares.AuthMiddleware, middlewares.RolesMiddleware(1), handlers.DirectoryHandler)
	startup.Init()
	gin.Run(":4001")
}
