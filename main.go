package main

import (
	"fileserver/handlers"
	"fileserver/middlewares"
	"fileserver/repository"
	"fileserver/startup"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}                                       // allowed hosts
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // allowed methods
	corsConfig.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"} // allowed headers
	corsConfig.AllowCredentials = true                                            // allow credentials

	repository.DBinitialize()
	gin := gin.Default()
	gin.Use(cors.New(corsConfig))
	gin.LoadHTMLGlob("templates/*")
	gin.Use(middlewares.DBMiddleware())
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
	gin.Run(startup.PORT)
}
