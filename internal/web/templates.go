package web

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/RiverMint78/pone-quest/ui"
)

// NewTemplateCache 预编译模板
func NewTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, fmt.Errorf("解析模板 %s 失败: %w", name, err)
		}
		cache[name] = ts
	}
	return cache, nil
}

// render 处理渲染逻辑和错误响应
func (h *Handler) render(w http.ResponseWriter, _ *http.Request, status int, page string, templateName string, data any) {
	ts, ok := h.TemplateCache[page]
	if !ok {
		http.Error(w, "模板不存在", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, "渲染失败", http.StatusInternalServerError)
	}
}
