//go:build !prod
// +build !prod

package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"jvanmelckebeke/anyconverter-api/docs"
)

func setupDocs(r *gin.Engine) {
	docs.SwaggerInfo.Title = "AnyConverter API"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Version = "1.0"

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// redirect /docs to /docs/index.html
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(301, "/docs/index.html")
	})
}
