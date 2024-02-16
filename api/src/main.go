package main

import (
	"github.com/gin-gonic/gin"
	"jvanmelckebeke/anyconverter-api/pkg/api"
	"jvanmelckebeke/anyconverter-api/pkg/constants"
	"jvanmelckebeke/anyconverter-api/pkg/logger"
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

	r.POST("/audio", api.PostAudio)

	if err := r.Run(":8000"); err != nil {
		logger.Error("error running gin", err)
		panic(err)
	}

}
