package templates

import (
	"embed"
)

//go:embed *.html
var PublicFs embed.FS
