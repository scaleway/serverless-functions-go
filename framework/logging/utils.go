package logging

import (
	"fmt"
	"strings"
)

func FormatLevel(i interface{}) string {
	if i == nil {
		i = "DEBUG"
	}
	return strings.ToUpper(fmt.Sprintf("%-5s -", i))
}

func FormatTimestamp(i interface{}) string {
	// we don't want to display timestamp in the log entry
	return ""
}
