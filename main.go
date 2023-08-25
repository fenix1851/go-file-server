package main

import (
	"fileserver/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	gin := gin.Default()
	gin.LoadHTMLGlob("templates/*")
	gin.GET("/", handlers.DirectoryHandler)
	gin.GET("/:path/*filepath", handlers.DirectoryHandler)
	gin.Run(":3001")
}
