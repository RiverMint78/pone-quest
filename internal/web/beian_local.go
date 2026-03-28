package web

import (
	"html/template"
	"os"
	"strings"
)

const beianLocalFile = "beian.local.html"

// LocalBeianHTML 读取本地备案片段；文件不存在时返回空字符串。
func LocalBeianHTML() template.HTML {
	raw, err := os.ReadFile(beianLocalFile)
	if err != nil {
		return ""
	}

	html := strings.TrimSpace(string(raw))
	if html == "" {
		return ""
	}

	return template.HTML(html)
}
