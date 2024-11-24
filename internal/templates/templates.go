package templates

import e "embed"

//go:embed all:data
var Templates e.FS
