package constants

import (
	"jvanmelckebeke/anyconverter-api/tools"
)

var UploadsDir = tools.GetEnv("UPLOADS_DIR", "/tmp")
