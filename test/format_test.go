package test

import (
	"testing"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/moorada/neferpitool/pkg/format"
)

func TestStringToTimeUTCOffsetSuffix(t *testing.T) {
	value := "2026-08-04 17:50:17 (UTC+8)"

	parsed, err := format.StringToTime(value)
	if err != nil {
		t.Fatalf("expected no parse error, got: %v", err)
	}

	_, offset := parsed.Zone()
	if offset != 8*60*60 {
		t.Fatalf("expected UTC+8 offset, got %d", offset)
	}

	if parsed.Year() != 2026 || parsed.Month() != 8 || parsed.Day() != 4 {
		t.Fatalf("unexpected parsed date: %v", parsed)
	}
}

func TestLimitDisplayWidthUnicode(t *testing.T) {
	value := "阿里云计算有限责任公司"

	truncated := format.LimitDisplayWidth(value, 12)
	if !utf8.ValidString(truncated) {
		t.Fatalf("expected valid utf-8 output, got: %q", truncated)
	}

	if runewidth.StringWidth(truncated) > 12 {
		t.Fatalf("expected width <= 12, got %d (%q)", runewidth.StringWidth(truncated), truncated)
	}
}
