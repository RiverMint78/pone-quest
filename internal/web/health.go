package web

import (
	"net/http"
	"text/template"
	"time"
)

var healthTmpl = template.Must(template.New("health").Parse(`[PoneQuest Status: OK]
Time: {{.Time}}

--- Request Info ---
RemoteAddr: {{.RemoteAddr}}
Method:     {{.Method}}
Host:       {{.Host}}
UA:         {{.UA}}
{{if .Forwarded}}Forwarded:  {{.Forwarded}}{{end}}`))

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	data := struct {
		Time       string
		RemoteAddr string
		Method     string
		Host       string
		UA         string
		Forwarded  string
	}{
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		RemoteAddr: r.RemoteAddr,
		Method:     r.Method,
		Host:       r.Host,
		UA:         r.UserAgent(),
		Forwarded:  r.Header.Get("X-Forwarded-For"),
	}

	_ = healthTmpl.Execute(w, data)
}
