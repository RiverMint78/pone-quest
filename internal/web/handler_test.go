package web

import (
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"
)

func BenchmarkHandler_Index(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := &Handler{
		Logger: logger,
	}

	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		h.handleIndex(w, req)
	}
}
