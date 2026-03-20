package web

import (
	"net/http"

	"github.com/RiverMint78/pone-quest/ui"
)

// RegisterRoutes 配置所有的路由映射
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	h.Logger.Debug("正在注册路由映射...")

	// 系统检查
	mux.HandleFunc("GET /healthz", h.handleHealth)

	// 业务逻辑
	bizMux := http.NewServeMux()
	bizMux.HandleFunc("GET /", h.handleIndex)
	bizMux.HandleFunc("GET /search", h.handleSearch)

	// 静态资源
	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("GET /static/", h.StaticCacheMiddleware(fileServer))

	// 图片资源
	imageDir := http.Dir("./data/images")
	imageServer := http.StripPrefix("/images/", http.FileServer(imageDir))
	mux.Handle("GET /images/", h.StaticCacheMiddleware(imageServer))

	// 添加 Recover
	mux.Handle("/", h.RecoverMiddleware(bizMux))

	h.Logger.Debug("路由注册完成")
}
