package main

import "strings"

func cropLongLabels(s string) string {
	if len(s) > 50 {
		return s[:47] + "..."
	}
	return s
}

func preprocessLabel(s string) string {
	if s == "" {
		return ""
	}
	s = cropLongLabels(s)
	if string(s[len(s)-1]) == "\\" {
		s += `\`
	}
	s = strings.Replace(s, `${`, `\${`, -1)
	s = strings.Replace(s, "`", "\\`", -1)
	return "`" + s + "`"
}
