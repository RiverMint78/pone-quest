package web

import (
	"net/http"

	"github.com/RiverMint78/pone-quest/ui"
)

// RegisterRoutes 配置所有的路由映射
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	h.Logger.Debug("正在注册路由映射...")

	// 静态资源
	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("GET /static/", fileServer)

	// 业务路由
	mux.HandleFunc("GET /", h.handleIndex)
	mux.HandleFunc("GET /search", h.handleSearch)

	h.Logger.Debug("路由注册完成")
}
