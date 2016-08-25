package main

func preprocess(i []string, o options) ([]string, options) {
	return i, o
}

type sectionType int
const (
	emptyT sectionType = iota
	stringT
	floatOrStringT
	separatorT
)

type lineFormat []sectionType

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
