package constants

import (
	"jvanmelckebeke/anyconverter-api/env"
)

var UploadsDir = env.GetEnv("UPLOADS_DIR", "/tmp")
