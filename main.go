package main

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	var (
		opts           = mustResolveOptions(os.Args[1:])
		rd   io.Reader = os.Stdin
	)
	if opts.lineFormat == "" {
		var newRd, lineFormat = readAndParseFormat(os.Stdin, opts.separator, opts.dateFormat)
		opts.lineFormat = lineFormat
		rd = newRd
	}
	dataset := mustNewDataset(rd, opts)
	opts.chartType = resolveChartType(opts.chartType, dataset.lf, dataset.fss, dataset.sss)
	if err := assertChartable(*dataset, opts); opts.debug || err != nil {
		showDebug(*dataset, opts, err)
		os.Exit(0)
	}
	b := newChartJSChart(*dataset, opts).mustBuild()
	tmpfile := mustNewTempFile()
	chartTempl := newChartTemplate(opts.chartType)
	chartTempl.mustExecute(b, tmpfile)
	tmpfile.mustClose()
	tmpfile.mustRenameWithHTMLSuffix()
	if err := open.Run(tmpfile.url()); err != nil {
		log.WithError(err).Fatalf("Could not open the default viewer; please configure open/xdg-open")
	}
}

func assertChartable(d dataset, opts options) error {
	switch opts.chartType {
	case pie:
		if len(d.fss) == 0 || (len(d.fss[0]) == 1 && len(d.sss) == 0 && len(d.tss) == 0) {
			return fmt.Errorf("couldn't find values to plot")
		}
	case bar:
		if len(d.fss) == 0 || (len(d.fss[0]) == 1 && len(d.sss) == 0 && len(d.tss) == 0) {
			return fmt.Errorf("couldn't find values to plot")
		}
	case line:
		if d.fss == nil || (d.sss == nil && d.tss == nil && len(d.fss[0]) < 2) {
			return fmt.Errorf("couldn't find values to plot")
		}
	case scatter:
		if len(d.fss) == 0 {
			return fmt.Errorf("couldn't find values to plot")
		}
	}
	return nil
}
