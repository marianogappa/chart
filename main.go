package main

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/marianogappa/chart/chartjs"
	"github.com/marianogappa/chart/format"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	var (
		opts           = mustResolveOptions(os.Args[1:])
		rd   io.Reader = os.Stdin
		err  error
	)
	if opts.rawLineFormat == "" {
		rd, opts.lineFormat = format.Parse(rd, opts.separator, opts.dateFormat)
	}
	dataset := mustNewDataset(rd, opts.lineFormat)
	if !opts.lineFormat.HasFloats && !opts.lineFormat.HasDateTimes && opts.lineFormat.HasStrings {
		dataset.fss, dataset.sss, opts.lineFormat = preprocessFreq(dataset.sss, opts.lineFormat)
	}
	if opts.chartType, err = resolveChartType(opts.chartType, opts.lineFormat, dataset.Len()); opts.debug || err != nil {
		fmt.Println(renderDebug(*dataset, opts, err))
		os.Exit(0)
	}
	b := chartjs.New(
		opts.chartType.String(),
		dataset.fss,
		dataset.sss,
		dataset.tss,
		dataset.minFSS,
		dataset.maxFSS,
		opts.title,
		opts.scaleType.String(),
		opts.xLabel,
		opts.yLabel,
		opts.zeroBased,
		int(opts.colorType),
	).MustBuild()
	tmpfile := mustNewTempFile()
	chartTempl := newChartTemplate(opts.chartType)
	chartTempl.mustExecute(b, tmpfile)
	tmpfile.mustClose()
	tmpfile.mustRenameWithHTMLSuffix()
	if err := open.Run(tmpfile.url()); err != nil {
		log.WithError(err).Fatalf("Could not open the default viewer; please configure open/xdg-open")
	}
}
