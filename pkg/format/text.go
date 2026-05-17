package format

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

func NormalizeText(s string) string {
	return strings.TrimSpace(strings.ToValidUTF8(s, ""))
}

func LimitDisplayWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}

	clean := NormalizeText(s)
	if runewidth.StringWidth(clean) <= width {
		return clean
	}

	tail := "..."
	if width <= runewidth.StringWidth(tail) {
		return runewidth.Truncate(clean, width, "")
	}

	return runewidth.Truncate(clean, width, tail)
}

func FitDisplayWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}

	return runewidth.FillRight(LimitDisplayWidth(s, width), width)
}
