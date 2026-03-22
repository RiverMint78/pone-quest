package pone

import (
	"strings"
	"unicode"
)

// NormalizeSearchText 统一规范化索引和查询文本：
// 1) 全角转半角
// 2) 引号统一
// 3) 合并多余空白
func NormalizeSearchText(input string) string {
	if input == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(input))

	for _, r := range input {
		r = normalizeWidthRune(r)
		r = normalizeQuoteRune(r)

		switch {
		case unicode.IsLetter(r), unicode.IsNumber(r):
			b.WriteRune(r)
		case unicode.IsSpace(r):
			b.WriteByte(' ')
		default:
			b.WriteRune(r)
		}
	}

	return strings.Join(strings.Fields(b.String()), " ")
}

func normalizeWidthRune(r rune) rune {
	// 全角空格
	if r == '\u3000' {
		return ' '
	}

	// U+FF01('！') ~ U+FF5E('～')
	if r >= '\uFF01' && r <= '\uFF5E' {
		return r - 0xFEE0
	}

	return r
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
