package main

func preprocess(i []string, o options) ([]string, options) {
	return i, o
}

func parseFormat(i []string, sep rune) string {
	lfs := make(map[string]int)
	for _, l := range i {
		lfs[parseLine(l, sep).string()] += 1
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

type sectionType int
const (
	emptyT sectionType = iota
	stringT
	floatOrStringT
	separatorT
)

type lineFormat []sectionType

func (lf lineFormat) string() string {
	s := ""
	for _, v := range lf {
		switch v {
		case stringT:
			s += "s"
		case floatOrStringT:
			s += "f"
		case separatorT:
			s += ","
		}
	}
	return s
}

func parseLine(s string, sep rune) lineFormat {
	lf := lineFormat{emptyT}
	for _, c := range s {
		switch lf[len(lf)-1] {
		case emptyT:
			if isFloatStart(c) {
				lf[len(lf)-1] = floatOrStringT
			} else if c == sep && sep != ' ' {
				lf[len(lf)-1] = floatOrStringT
				lf = append(lf, separatorT)
			} else if !(c == sep) {
				lf[len(lf)-1] = stringT
			}
		case stringT:
			if c == sep {
				lf = append(lf, separatorT)
			}
		case floatOrStringT:
			if c == sep {
				lf = append(lf, separatorT)
			} else if !isFloat(c) && !(c == sep) {
				lf[len(lf)-1] = stringT
			}
		case separatorT:
			if isFloatStart(c) {
				lf = append(lf, floatOrStringT)
			} else if c == sep && sep != ' ' {
				lf = append(lf, floatOrStringT, separatorT)
			} else if sep != ' ' {
				lf = append(lf, stringT)
			}
		}
	}
	if sep == ' ' && lf[len(lf)-1] == separatorT {
		return lf[:len(lf)-1]
	} else if lf[len(lf)-1] == separatorT {
		return append(lf, floatOrStringT)
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
