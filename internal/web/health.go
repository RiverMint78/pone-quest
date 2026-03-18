package web

import (
	"net/http"
	"sort"
	"strings"
	"text/template"
	"time"
)

var healthTmpl = template.Must(template.New("health").Parse(`
[PoneQuest System Status: {{.Status}}]
Timestamp:  {{.Time}}
Uptime:     {{.Uptime}}

--- Client Header ---
{{range .Headers}}{{.Key}}: {{.Value}}
{{end}}
--- Summary ---
Total-Headers: {{.HeaderCount}}
Content-Length: {{.ContentLength}}
TLS-Version:    {{.TLS}}
`))

var startTime = time.Now()

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	type headerKV struct{ Key, Value string }
	var headers []headerKV
	for k, v := range r.Header {
		lowerK := strings.ToLower(k)
		if strings.Contains(lowerK, "auth") || strings.Contains(lowerK, "cookie") || strings.Contains(lowerK, "token") {
			headers = append(headers, headerKV{k, "****** (Redacted)"})
			continue
		}
		headers = append(headers, headerKV{k, strings.Join(v, ", ")})
	}
	sort.Slice(headers, func(i, j int) bool { return headers[i].Key < headers[j].Key })

	tlsVer := "None (Plain HTTP)"
	if r.TLS != nil {
		tlsVer = "Encrypted (HTTPS)"
	}

	data := struct {
		Status        string
		Time          string
		Uptime        string
		Headers       []headerKV
		HeaderCount   int
		ContentLength int64
		TLS           string
	}{
		Status:        "OPERATIONAL",
		Time:          time.Now().Format(time.RFC3339),
		Uptime:        time.Since(startTime).Truncate(time.Second).String(),
		Headers:       headers,
		HeaderCount:   len(r.Header),
		ContentLength: r.ContentLength,
		TLS:           tlsVer,
	}

	_ = healthTmpl.Execute(w, data)
}
