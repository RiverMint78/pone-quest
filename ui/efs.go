// Package ui: embed FS
package ui

import "embed"

//go:embed "html" "static"
var Files embed.FS
