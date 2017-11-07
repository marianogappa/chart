package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func readAndParseFormat(r io.Reader, sep rune, df string) ([]string, string) {
	var (
		ls  = make([]string, 0, 500)
		lfs = make(map[string]int)
		rd  = bufio.NewReader(r)
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
		ls = append(ls, l)
		lfs[parseLineFormat(l, sep, df)]++
	}
	return ls, maxLineFormat(lfs)
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
	var sp []string
	s = string(regexp.MustCompile(string(sep)+"{2,}").ReplaceAll([]byte(s), []byte(string(sep))))
	sp = strings.Split(strings.TrimSpace(s), string(sep))
	// In certain cases data may appear to be tabbed, but
	// instead there is an irregular number of spaces used
	// as separator, rather than a tabstop. In such a case
	// parsing above may result in a slice with a single
	// element, not correctly parsed. In this case we want
	// to attempt split using strings.Fields as opposed to
	// strings.Split.
	if len(sp) == 1 && strings.Count(sp[0], " ") > 0 {
		sp = strings.Fields(strings.TrimSpace(s))
	}
	lf := ""
	for _, sc := range sp {
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

func parseLine(l string, lf string, sep rune, df string) ([]float64, []string, []time.Time, error) {
	var sp []string
	l = string(regexp.MustCompile(string(sep)+"{2,}").ReplaceAll([]byte(l), []byte(string(sep))))
	sp = strings.Split(strings.TrimSpace(l), string(sep))
	// See explanation in parseLineFormat for this seemingly
	// unnecessary splitting code.
	if len(sp) == 1 && strings.Count(sp[0], " ") > 0 {
		sp = strings.Fields(strings.TrimSpace(l))
	}

	fs := []float64{}
	ss := []string{}
	ds := []time.Time{}

	if len(sp) < len(lf) {
		return fs, ss, ds, fmt.Errorf("Input line has invalid format length; expected %v vs found %v", len(lf), len(sp))
	}

	for i, lfe := range lf {

		s := strings.TrimSpace(sp[i])
		switch lfe {
		case 's':
			ss = append(ss, s)
		case 'd':
			d, err := time.Parse(df, s)
			if err != nil {
				return fs, ss, ds, fmt.Errorf("Couldn't convert %v to date given: %v", s, err)
			}
			ds = append(ds, d)
		case 'f':
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return fs, ss, ds, fmt.Errorf("Couldn't convert %v to float given: %v", s, err)
			}
			fs = append(fs, f)
		}
	}

	return fs, ss, ds, nil
}
