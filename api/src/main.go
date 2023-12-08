package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jvanmelckebeke/anyconverter-api/api"
	"jvanmelckebeke/anyconverter-api/constants"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "jvanmelckebeke/anyconverter-api/docs"
)

func main() {
	r := gin.Default()

	docs.SwaggerInfo.Title = "AnyConverter API"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Version = "1.0"

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// redirect /docs to /docs/index.html
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(301, "/docs/index.html")
	})

	r.GET(constants.AllStatusEndpoint,
		api.GetTasks)

	r.GET(constants.SingleStatusEndpoint, api.GetTaskStatus)

	r.GET(
		constants.SingleResultEndpoint,
		api.GetResult)

	r.POST("/image", api.PostImage)

	r.POST("/video", api.PostVideo)

	if err := r.Run(":8000"); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
