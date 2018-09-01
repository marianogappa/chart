package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/marianogappa/chart/chartjs"
	"github.com/marianogappa/chart/dataset"
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
	dataset := dataset.MustNew(rd, opts.lineFormat)
	if !opts.lineFormat.HasFloats && !opts.lineFormat.HasDateTimes && opts.lineFormat.HasStrings {
		dataset.FSS, dataset.SSS, opts.lineFormat = preprocessFreq(dataset.SSS, opts.lineFormat)
	}
	if opts.chartType, err = resolveChartType(opts.chartType, opts.lineFormat, dataset.Len()); opts.debug || err != nil {
		fmt.Println(renderDebug(*dataset, opts, err))
		os.Exit(0)
	}
	chart := chartjs.New(
		chartjs.NewChartType(opts.chartType.String()),
		*dataset,
		chartjs.Options{
			Title:     opts.title,
			ScaleType: chartjs.NewScaleType(opts.scaleType.String()),
			XLabel:    opts.xLabel,
			YLabel:    opts.yLabel,
			ZeroBased: opts.zeroBased,
			ColorType: chartjs.NewColorType(opts.colorType.String()),
		},
	)
	tmpfile := mustNewTempFile()
	chart.MustBuild(chartjs.OutputAll, tmpfile.f)
	tmpfile.mustClose()
	tmpfile.mustRenameWithHTMLSuffix()
	if err := open.Run(tmpfile.url()); err != nil {
		log.Fatalf("Could not open the default viewer; please configure open/xdg-open: %v", err)
	}
}
