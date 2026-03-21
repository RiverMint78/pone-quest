package pone

import (
	"strings"
	"unicode"
)

// NormalizeSearchText 统一规范化索引和查询文本：
// 1) 全部转小写
// 2) 删除特殊符号（保留 . , ? ! - 和引号）
// 3) 各类引号统一为标准双引号 (")
// 4) 合并多余空白
func NormalizeSearchText(input string) string {
	if input == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(input))

	for _, r := range input {
		r = normalizeQuoteRune(r)

		switch {
		case unicode.IsLetter(r), unicode.IsNumber(r):
			b.WriteRune(r)
		case unicode.IsSpace(r):
			b.WriteByte(' ')
		case r == '.', r == ',', r == '?', r == '!', r == '-', r == '"', r == '\'':
			b.WriteRune(r)
		}
	}

	out := strings.ToLower(b.String())
	return strings.Join(strings.Fields(out), " ")
}

func normalizeQuoteRune(r rune) rune {
	switch r {
	case '"', '`',
		'“', '”', '„', '‟',
		'‘', '’', '‚', '‛',
		'«', '»',
		'「', '」', '『', '』',
		'＂', '＇':
		return '"'
	default:
		return r
	}
}
