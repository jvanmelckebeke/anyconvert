package constants

import (
	"jvanmelckebeke/anyconverter-api/pkg/env"
)

var UploadsDir = env.Getenv("UPLOADS_DIR", "/tmp")

var OutputDir = env.Getenv("OUTPUT_DIR", "/tmp")
