package main

func preprocess(i []string, o options) ([]string, options) {
	return i, o
}

func parseFormat(i []string, sep rune) string {
	lfs := make(map[string]int)
	for _, l := range i {
		lfs[parseLine(l, sep)] += 1
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

func parseLine(s string, sep rune) string {
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
			} else if sep != ' ' {
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
