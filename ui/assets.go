package ui

import (
	"crypto/sha256"
	"encoding/hex"
	"path"
	"strings"
	"sync"
)

var (
	assetHashOnce sync.Once
	assetHashMap  map[string]string
)

func buildAssetHashes() {
	assetHashMap = map[string]string{}

	entries := []string{
		"static/css/style.css",
		"static/js/ui.js",
		"static/favicon.svg",
	}

	for _, name := range entries {
		b, err := Files.ReadFile(name)
		if err != nil {
			continue
		}
		sum := sha256.Sum256(b)
		assetHashMap[name] = hex.EncodeToString(sum[:])[:12]
	}
}

// AssetPath returns cache-busted static URLs like /static/css/style.css?v=<hash>.
func AssetPath(rel string) string {
	assetHashOnce.Do(buildAssetHashes)

	clean := strings.TrimLeft(path.Clean(rel), "/")
	if after, ok := strings.CutPrefix(clean, "static/"); ok {
		clean = after
	}

	full := "static/" + clean
	url := "/" + full

	if h, ok := assetHashMap[full]; ok {
		return url + "?v=" + h
	}

	return url
}
