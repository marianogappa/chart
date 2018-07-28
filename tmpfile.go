package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

type tempFile struct {
	f    *os.File
	name string
}

func mustNewTempFile() *tempFile {
	tmpfile, err := ioutil.TempFile("", "chartData")
	if err != nil {
		log.WithField("err", err).Fatalf("Could not create temporary file to store the chart.")
	}
	return &tempFile{tmpfile, tmpfile.Name()}
}

func (f *tempFile) mustClose() {
	if err := f.f.Close(); err != nil {
		log.WithField("err", err).Fatalf("Could not close temporary file after saving chart to it.")
	}
}

func (f *tempFile) mustRenameWithHTMLSuffix() {
	newName := f.name + ".html"
	if err := os.Rename(f.name, newName); err != nil {
		log.WithField("err", err).Fatalf("Could not add html extension to the temporary file.")
	}
	f.name = newName
}

func (f *tempFile) url() string {
	return fmt.Sprintf("file://%v", f.name)
}
