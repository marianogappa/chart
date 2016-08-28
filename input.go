package main

import (
	"bufio"
	"io"
	"strings"
)

func readInput(r io.Reader) []string {
	ls := []string{}
	var err error
	rd := bufio.NewReader(r)
	for {
		var s string
		s, err = rd.ReadString('\n')
		if err == io.EOF {
			return ls
		}
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		ls = append(ls, strings.TrimSpace(s))
	}
}
