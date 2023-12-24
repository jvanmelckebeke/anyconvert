package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jvanmelckebeke/anyconverter-api/pkg/api"
	"jvanmelckebeke/anyconverter-api/pkg/constants"
)

func main() {
	r := gin.Default()

	setupDocs(r)

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
