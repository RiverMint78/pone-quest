package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"sync"
	"time"

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
		slog.Debug("已编译模板页面", slog.String("name", name), slog.Any("patterns", patterns))
	}
	slog.Debug("模板缓存初始化完成", slog.Int("total_pages", len(cache)))
	return cache, nil
}

var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

// render 处理渲染逻辑和错误响应
func (h *Handler) render(w http.ResponseWriter, r *http.Request, status int, page string, templateName string, data any) {
	start := time.Now()

	ts, ok := h.TemplateCache[page]
	if !ok {
		http.Error(w, "模板不存在", http.StatusInternalServerError)
		return
	}

	// Buffer from Pool
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	err := ts.ExecuteTemplate(buf, templateName, data)

	if err != nil {
		h.Logger.ErrorContext(r.Context(), "执行模板渲染失败",
			slog.String("page", page),
			slog.String("target", templateName),
			slog.Any("err", err))

		http.Error(w, "服务器内部渲染错误", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)

	h.Logger.DebugContext(r.Context(), "模板渲染成功",
		slog.String("page", page),
		slog.String("target", templateName),
		slog.Duration("latency", time.Since(start)),
	)
}
