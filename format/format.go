// Package format is used for inferring the line format
// of a dataset with one line per data point.
package format

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Parse infers the line format of a dataset with one line per data point.
// It requires a separator e.g. `\t` and optionally a date format that time.Parse
// understands. Parse also returns a new io.Reader ready to consume again.
func Parse(r io.Reader, separator rune, dateFormat string) (io.Reader, string) {
	var (
		lfs = make(map[string]int)
		rd  = bufio.NewReader(r)
		buf bytes.Buffer
		l   string
		err error
	)
	for {
		l, err = rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		buf.WriteString(l)
		buf.WriteByte('\n')
		lfs[parseLineFormat(l, separator, dateFormat)]++
	}
	return &buf, maxLineFormat(lfs)
}

func maxLineFormat(lfs map[string]int) string {
	max := 0
	lf := ""
	for k, v := range lfs {
		if v > max {
			max = v
			lf = k
		}
	}
	return lf
}

func parseLineFormat(s string, sep rune, df string) string {
	s = string(regexp.MustCompile(string(sep)+"{2,}").ReplaceAll([]byte(s), []byte(string(sep))))
	ss := strings.Split(strings.TrimSpace(s), string(sep))
	lf := ""
	for _, sc := range ss {
		if _, err := strconv.ParseFloat(sc, 64); err == nil {
			lf += "f"
		} else if _, err := time.Parse(df, sc); err == nil && strings.TrimSpace(sc) != "" {
			lf += "d"
		} else {
			lf += "s"
		}
	}
	return lf
}
