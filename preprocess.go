package main

import (
	"fmt"
	"strconv"
	"strings"
)

func preprocess(i []string, o options) ([][]float64, [][]string, options) {
	var fss [][]float64
	var sss [][]string

	sep := rune(o.separator[0])
	lf := parseFormat(i, sep)
	for _, l := range i {
		fs, ss, err := parseLine(l, lf, sep)
		if err != nil {
			break
		}
		fss = append(fss, fs)
		sss = append(sss, ss)
	}

	if strings.Index(lf, "f") == -1 {
		fss, sss = preprocessFreq(sss)
	}

	return fss, sss, o
}

func parseFormat(i []string, sep rune) string {
	lfs := make(map[string]int)
	for _, l := range i {
		if len(strings.TrimSpace(l)) == 0 {
			continue
		}
		lfs[parseLineFormat(l, sep)] += 1
	}
	return maxLineFormat(lfs)
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

func parseLineFormat(s string, sep rune) string {
	lf := " "
	for _, c := range s {
		switch lf[len(lf)-1] {
		case ' ':
			if isFloatStart(c) {
				lf = "f"
			} else if c == sep && sep != ' ' {
				lf = "f,"
			} else if !(c == sep) {
				lf = "s"
			}
		case 's':
			if c == sep {
				lf += ","
			}
		case 'f':
			if c == sep {
				lf += ","
			} else if !isFloat(c) && !(c == sep) {
				lf = lf[:len(lf)-1] + "s"
			}
		case ',':
			if isFloatStart(c) {
				lf += "f"
			} else if c == sep && sep != ' ' {
				lf += "f,"
			} else if c != sep {
				lf += "s"
			}
		}
	}
	if sep == ' ' && lf[len(lf)-1] == ',' {
		return lf[:len(lf)-1]
	} else if lf[len(lf)-1] == ',' {
		return lf + "f"
	}
	return lf
}

func isFloat(c rune) bool {
	if c == '.' || c == 'e' || c == 'E' || c == '-' || c == '0' || c == '1' || c == '2' || c == '3' ||
		c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9' {
		return true
	}
	return false
}
func isFloatStart(c rune) bool {
	if c == '-' || c == '0' || c == '1' || c == '2' || c == '3' ||
		c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9' {
		return true
	}
	return false
}

func parseLine(l string, lf string, sep rune) ([]float64, []string, error) {
	var pf string
	var ps string
	fs := []float64{}
	ss := []string{}
	lfi := 0
	li := 0
	for {
		if li > len(l)-1 {
			break
		}
		lv := rune(l[li])

		if lfi > len(lf)-1 {
			return fs, ss, nil
		}
		switch lf[lfi] {
		case 'f':
			if isFloat(lv) || (lv == ' ' && sep != ' ') {
				pf += string(lv)
			}
			if !(isFloat(lv) || (lv == ' ' && sep != ' ')) || li == len(l)-1 {
				if len(pf) == 0 {
					return fs, ss, fmt.Errorf("Parse error; expecting float but got [%v]", string(lv))
				}
				nf, err := strconv.ParseFloat(strings.TrimSpace(pf), 64)
				if err != nil {
					return fs, ss, fmt.Errorf("Parse error; invalid float [%v]", string(pf))
				}
				fs = append(fs, nf)
				pf = ""
				lfi++
				li--
			}
		case 's':
			if lv != sep {
				ps += string(lv)
			}
			if lv == sep || li == len(l)-1 {
				if len(ps) == 0 {
					return fs, ss, fmt.Errorf("Parse error; expecting string but got [%v]", string(lv))
				}
				ss = append(ss, escapeString(ps))
				ps = ""
				lfi++
				li--
			}
		case ',':
			if lv == sep {
				lfi++
			} else {
				return fs, ss, fmt.Errorf("Parse error; expecting [%v], but found [%v]", string(sep), string(lv))
			}
		}
		li++
	}
	if lfi < len(lf)-1 {
		return fs, ss, fmt.Errorf("Parse error; unfinished line [%v] according to format [%v]", l, lf)
	}
	return fs, ss, nil
}

func escapeString(s string) string {
	if string(s[len(s)-1]) == "\\" {
		s += `\`
	}
	s = strings.TrimSpace(s)
	s = strings.Replace(s, `${`, `\${`, -1)
	s = strings.Replace(s, "`", "\\`", -1)
	return s
}
