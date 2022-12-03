// Package dataset contains the Dataset struct which represents the source dataset used to construct the charts.
// Use it to read a dataset from an io.Reader, with the New function.
package dataset

import (
	"bufio"
	"io"
	"log"
	"time"

	"github.com/marianogappa/chart/format"
)

// Dataset represents the source dataset used to construct the charts.
// Use it to read a dataset from an io.Reader, with the New function.
type Dataset struct {
	FSS      [][]float64
	SSS      [][]string
	TSS      [][]time.Time
	StdinLen int
}

// Len returns the number of datapoints in the Dataset. Strings alone
// are not considered datapoints.
func (d Dataset) Len() int {
	if d.FSS == nil {
		return len(d.TSS)
	}
	return len(d.FSS)
}

// MustNew reads an io.Reader expected to have a dataset with one
// datapoint per line and returns the parsed Dataset. Fatals on
// error. In order to parse each line, it requires a LineFormat.
//
// For example, for an input like this:
//
// ABC,1,2,2001-02-03
// DEF,3,4,2004-05-06
// ...
//
// This invocation would be appropriate:
//
// ds := dataset.MustNew(format.Parse(os.Stdin, ",", "2006-01-02"))
func MustNew(r io.Reader, f format.LineFormat) *Dataset {
	d, err := New(r, f)
	if err != nil {
		log.Fatalf("Could not build dataset: %v", err)
	}
	return d
}

// New reads an io.Reader expected to have a dataset with one
// datapoint per line and returns the parsed Dataset. In order
// to parse each line, it requires a LineFormat.
//
// For example, for an input like this:
//
// ABC,1,2,2001-02-03
// DEF,3,4,2004-05-06
// ...
//
// This invocation would be appropriate:
//
// ds, err := dataset.New(format.Parse(os.Stdin, ",", "2006-01-02"))
func New(r io.Reader, f format.LineFormat) (*Dataset, error) {
	d := &Dataset{
		FSS: make([][]float64, 0, 500),
		SSS: make([][]string, 0, 500),
		TSS: make([][]time.Time, 0, 500),
	}
	return d, d.read(r, f)
}

func (d *Dataset) read(r io.Reader, f format.LineFormat) error {
	var (
		nilSSS, nilFSS, nilTSS = true, true, true
		scanner                = bufio.NewScanner(r)
	)
	for scanner.Scan() {
		d.StdinLen++
		fs, ss, ts, err := f.ParseLine(scanner.Text())
		if err != nil {
			continue
		}
		if nilSSS && len(ss) > 0 {
			nilSSS = false
		}
		if nilFSS && len(fs) > 0 {
			nilFSS = false
		}
		if nilTSS && len(ts) > 0 {
			nilTSS = false
		}
		d.FSS = append(d.FSS, fs)
		d.SSS = append(d.SSS, ss)
		d.TSS = append(d.TSS, ts)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if nilSSS {
		d.SSS = nil
	}
	if nilFSS {
		d.FSS = nil
	}
	if nilTSS {
		d.TSS = nil
	}

	return nil
}
