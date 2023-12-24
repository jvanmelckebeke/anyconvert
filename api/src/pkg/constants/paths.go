package constants

import (
	"jvanmelckebeke/anyconverter-api/pkg/env"
)

var UploadsDir = env.GetEnv("UPLOADS_DIR", "/tmp")

var OutputDir = env.GetEnv("OUTPUT_DIR", "/tmp")
