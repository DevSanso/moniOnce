package funcs

import (
	"strconv"
	"strings"
)

func ParseBytesFromStr(s string) int64 {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "KiB") {
		f, _ := strconv.ParseFloat(strings.TrimSuffix(s, "KiB"), 64)
		return int64(f * 1024)
	}
	if strings.HasSuffix(s, "MiB") {
		f, _ := strconv.ParseFloat(strings.TrimSuffix(s, "MiB"), 64)
		return int64(f * 1024 * 1024)
	}
	if strings.HasSuffix(s, "bytes") {
		f, _ := strconv.ParseFloat(strings.TrimSuffix(s, "bytes"), 64)
		return int64(f)
	}
	return 0
}

func ParseBool(val string) bool {
	val = strings.TrimSpace(strings.ToLower(val))
	return val == "true" || val == "yes" || val == "1"
}