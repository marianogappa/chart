package main

import (
	"io"
	"bufio"
	"log"
)

func mustReadInput(r io.Reader) []string {
	ls, err := readInput(r)
	if err != nil {
		log.Fatal(err)
	}
	return ls
}

func readInput(r io.Reader) ([]string, error) {
	var ls []string
	var err error
	rd := bufio.NewReader(r)
	for {
		var s string
		s, err = rd.ReadString('\n')
		if err == io.EOF {
			return ls, nil
		}
		if err != nil {
			return ls, err
		}
		ls = append(ls, s)
	}
}
