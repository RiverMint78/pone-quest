// Package ui: embed FS
package ui

import "embed"

//go:embed "html" "static/css" "static/js"
var Files embed.FS
