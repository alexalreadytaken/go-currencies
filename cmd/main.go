package main

import (
	_ "github.com/alexalreadytaken/go-currencies/docs"
	ginpkg "github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title Currencies service api
// @version 1.0
// @BasePath /api
// @schemes http
func main() {
	println("hello world")
	gin := ginpkg.Default()
	api := gin.Group("/api")
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // /api/docs/index.html
	gin.Run(":2000")
}
