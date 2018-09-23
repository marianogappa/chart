package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type tempFile struct {
	f    *os.File
	name string
}

func mustNewTempFile() *tempFile {
	tmpfile, err := ioutil.TempFile("", "chartData")
	if err != nil {
		log.Fatalf("Could not create temporary file to store the chart: %v", err)
	}
	return &tempFile{tmpfile, tmpfile.Name()}
}

func (f *tempFile) mustClose() {
	if err := f.f.Close(); err != nil {
		log.Fatalf("Could not close temporary file after saving chart to it: %v", err)
	}
}

func (f *tempFile) mustRenameWithHTMLSuffix() {
	newName := f.name + ".html"
	if err := os.Rename(f.name, newName); err != nil {
		log.Fatalf("Could not add html extension to the temporary file: %v", err)
	}
	f.name = newName
}

func (f *tempFile) url() string {
	return fmt.Sprintf("file://%v", f.name)
}
