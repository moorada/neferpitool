package format

import (
	"fmt"
	"github.com/araddon/dateparse"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*Return the time in format like whois standard*/
func TimeNowStringWhoisFormat() string {
	t := time.Now().UTC()
	return t.Format("2006-01-02T15:04:05Z07:00")
}

func TimeToStringConsole(input time.Time) string {
	return input.Format("02/01/2006")
}

/*Convert time string in whois stadard to time*/
func StringToTime(d string) (time.Time, error) {
	d = normalizeDateString(d)

	layouts := []string{
		"02-Jan-2006",
		"2006-01-02",
		"02-01-2006",
		"Jan-02-2006",
		"2006-01-02",
		"02-01-2006 15:04:05",
		"2006-01-02 15:04:05",
		"January-2-2006",
		"January--2-2006",
		"2-Jan-2006",
		"2006-1-2",
		"2-1-2006",
		"Jan-2-2006",
		"2006-1-2",
		"2-1-2006 15:04:05",
		"January--02-2006",
		"January-02-2006",
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05 MST",
	}

	separators := []string{
		".",
		",",
		"/",
		":",
		"-",
		" ",
	}

	var layoutsTotal []string

	for _, l := range layouts {
		for _, s := range separators {
			layoutsTotal = append(layoutsTotal, strings.Replace(l, "-", s, -1))
		}
	}

	//special cases
	layoutsTotal = append(layoutsTotal, "2006. 01. 02.")
	//

	t, err := dateparse.ParseAny(d)
	if err != nil {
		err = nil
		for _, l := range layoutsTotal {
			t, err = time.Parse(l, d)
			if err == nil {
				break
			}
		}
	}

	//t, err := dateparse.ParseAny(timeWhoisFormat)

	return t, err
}

var utcOffsetSuffixPattern = regexp.MustCompile(`\s*\(UTC([+-]\d{1,2})(?::?(\d{2}))?\)\s*$`)

func normalizeDateString(input string) string {
	value := strings.TrimSpace(strings.ToValidUTF8(input, ""))

	matches := utcOffsetSuffixPattern.FindStringSubmatch(value)
	if len(matches) == 0 {
		return value
	}

	hour, err := strconv.Atoi(matches[1])
	if err != nil {
		return value
	}

	minutes := 0
	if matches[2] != "" {
		minutes, err = strconv.Atoi(matches[2])
		if err != nil {
			minutes = 0
		}
	}

	offset := fmt.Sprintf("%+03d:%02d", hour, minutes)
	return strings.TrimSpace(utcOffsetSuffixPattern.ReplaceAllString(value, " "+offset))
}
