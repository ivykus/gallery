package templates

import (
	"embed"
	_ "embed"
)

//go:embed *.gohtml
var FS embed.FS
